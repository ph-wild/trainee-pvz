package integration

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"trainee-pvz/config"
	"trainee-pvz/internal/handler"
	"trainee-pvz/internal/repository"
	"trainee-pvz/internal/service"
)

var cities = []string{"Москва", "Санкт-Петербург", "Казань"}
var productTypes = []string{"одежда", "обувь", "электроника"}

type fakeMetrics struct{}

func (f *fakeMetrics) SaveEntityCount(value float64, entity string) {}

func (f *fakeMetrics) SaveHTTPDuration(timeSince time.Time, path string, code int, method string) {}

func randomCity(r *rand.Rand) string {
	return cities[r.IntN(len(cities))]
}

func randomProductType(r *rand.Rand) string {
	return productTypes[r.IntN(len(productTypes))]
}
func TestEndToEndPVZFlow(t *testing.T) {
	cfg, err := config.GetConfig("../../config.yaml")
	require.NoError(t, err)

	db, err := sqlx.Connect("postgres", cfg.DB.Connection)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	pvzRepo := repository.NewPVZRepository(db)
	receptionRepo := repository.NewReceptionRepository(db)
	productRepo := repository.NewProductRepository(db)

	userService := service.NewUserService(userRepo)
	pvzService := service.NewPVZService(pvzRepo, &fakeMetrics{})
	receptionService := service.NewReceptionService(receptionRepo, &fakeMetrics{})
	productService := service.NewProductService(productRepo, &fakeMetrics{})

	services := handler.Services{
		User:      userService,
		PVZ:       pvzService,
		Reception: receptionService,
		Product:   productService,
	}

	s := handler.NewServer(services, nil, cfg, &fakeMetrics{}) //  without JWT

	srv := httptest.NewServer(s.Routes())
	defer srv.Close()

	// 1. CREATE PVZ
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	city := randomCity(r)
	pvzID := uuid.New()
	pvzPayload := map[string]any{
		"id":               pvzID.String(),
		"city":             city,
		"registrationDate": time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(pvzPayload)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/pvz", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer dummy-token-for-moderator")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created struct {
		Id string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&created)
	pvzID, err = uuid.Parse(created.Id)
	require.NoError(t, err)

	// 2. CREATE Reception
	recBody := map[string]string{"pvzId": pvzID.String()}
	body, _ = json.Marshal(recBody)
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/receptions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer dummy-token-for-employee")
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// 3. Add 50 Items
	for i := 0; i < 50; i++ {
		productType := randomProductType(r)
		product := map[string]any{
			"pvzId": pvzID.String(),
			"type":  productType,
		}
		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/products", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer dummy-token-for-employee")
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	// 4. Close Reception
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/pvz/"+pvzID.String()+"/close_last_reception", nil)
	req.Header.Set("Authorization", "Bearer dummy-token-for-employee")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// 5. Clean Data
	_, err = db.Exec(`DELETE FROM products WHERE reception_id IN (SELECT id FROM receptions WHERE pvz_id = $1)`, pvzID)
	require.NoError(t, err)
	_, err = db.Exec(`DELETE FROM receptions WHERE pvz_id = $1`, pvzID)
	require.NoError(t, err)
	_, err = db.Exec(`DELETE FROM pvz WHERE id = $1`, pvzID)
	require.NoError(t, err)
}
