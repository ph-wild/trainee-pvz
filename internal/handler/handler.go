package handler

import (
	"context"
	"encoding/json"

	"log/slog"
	"net/http"
	"strconv"
	"time"

	"trainee-pvz/api"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"trainee-pvz/config"
	"trainee-pvz/internal/auth"
	er "trainee-pvz/internal/errors"
	"trainee-pvz/internal/models"
	"trainee-pvz/internal/openapi"
)

type UserServiceInterface interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, email string) (models.User, error)
}

type ProductServiceInterface interface {
	AddProduct(ctx context.Context, p models.Product) error
	DeleteLastProduct(ctx context.Context, receptionID string) error
}

type ReceptionServiceInterface interface {
	CreateReception(ctx context.Context, rec models.Reception) error
	CloseReception(ctx context.Context, id string) error
	GetLastReceptionID(ctx context.Context, pvzID string) (string, error)
	GetOpenReceptionID(ctx context.Context, pvzID string) (string, error)
}

type PVZServiceInterface interface {
	CreatePVZ(ctx context.Context, pvz models.PVZ) error
	ListPVZ(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
}

type Server struct {
	Service    Services
	JWTManager *auth.JWTManager
	Cfg        config.Cfg
	metrics    metrics
}

type Services struct {
	User      UserServiceInterface
	Product   ProductServiceInterface
	Reception ReceptionServiceInterface
	PVZ       PVZServiceInterface
}

type metrics interface {
	SaveHTTPDuration(timeSince time.Time, path string, code int, method string)
}

func NewServer(service Services, jwt *auth.JWTManager, cfg config.Cfg, m metrics) *Server {
	return &Server{
		Service:    service,
		JWTManager: jwt,
		Cfg:        cfg,
		metrics:    m,
	}
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostRegisterJSONBody
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("failed to decode register request", slog.Any("err", err))
		http.Error(w, `{"message":"invalid json"}`, http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:       uuid.New().String(),
		Email:    string(req.Email),
		Password: string(hashed),
		Role:     string(req.Role),
	}

	err = s.Service.User.Register(ctx, user)
	if err != nil {
		if errors.Is(err, er.ErrUserAlreadyExists) {
			http.Error(w, `{"message":"user already exists"}`, http.StatusBadRequest)
			return
		}
		slog.Error("failed to create user", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	token, err := s.JWTManager.Generate(user.ID, user.Role)
	if err != nil {
		slog.Error("failed to generate jwt token", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostLoginJSONBody

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("failed to decode login request", slog.Any("err", err))
		http.Error(w, `{"message":"invalid json"}`, http.StatusBadRequest)
		return
	}

	user, err := s.Service.User.Login(ctx, string(req.Email))
	if err != nil {
		slog.Error("user not found", slog.Any("err", err))
		http.Error(w, `{"message":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		slog.Warn("password mismatch", slog.String("email", user.Email))
		http.Error(w, `{"message":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := s.JWTManager.Generate(user.ID, user.Role)
	if err != nil {
		slog.Error("failed to generate token", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostDummyLoginJSONBody

	_, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("invalid dummy login body", slog.Any("err", err))
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	token := s.Cfg.Auth.DummyTokenPrefix + string(req.Role)

	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostPvzJSONRequestBody

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("invalid pvz json", slog.Any("err", err))
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		//metrics.SaveHTTPCount(1, r.URL.Path, http.StatusBadRequest, r.Method)
		return
	}

	city := string(req.City)
	id := uuid.New()
	now := time.Now().UTC()

	pvz := models.PVZ{
		ID:               id.String(),
		RegistrationDate: now,
		City:             city,
	}

	err = s.Service.PVZ.CreatePVZ(ctx, pvz)
	if errors.Is(err, er.ErrUnsupportedCity) {
		http.Error(w, `{"message":"unsupported city"}`, http.StatusBadRequest)
		//metrics.SaveHTTPCount(1, r.URL.Path, http.StatusInternalServerError, r.Method)
		return
	}

	if err != nil {
		slog.Error("failed to create pvz", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("pvz has been created", slog.Any("info:", pvz))

	openapiUUID := openapi_types.UUID(id)
	resp := openapi.PVZ{
		Id:               &openapiUUID,
		RegistrationDate: &now,
		City:             openapi.PVZCity(city),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) CreateReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostReceptionsJSONRequestBody

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("invalid reception json", slog.Any("err", err))
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	pvzID := req.PvzId.String()
	now := time.Now().UTC()
	id := uuid.New()

	reception := models.Reception{
		ID:       id.String(),
		DateTime: now,
		PVZID:    pvzID,
		Status:   string(openapi.InProgress),
	}

	err = s.Service.Reception.CreateReception(ctx, reception)
	if errors.Is(err, er.ErrReceptionAlreadyExists) {
		http.Error(w, `{"message":"can't create one more reception"}`, http.StatusBadRequest)
		return
	}

	if err != nil {
		slog.Error("failed to create reception", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}
	slog.Info("reception has been created", slog.Any("info:", reception))

	openapiID := openapi_types.UUID(id)
	resp := openapi.Reception{
		Id:       &openapiID,
		DateTime: now,
		PvzId:    req.PvzId,
		Status:   openapi.InProgress,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostProductsJSONRequestBody

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("invalid product json", slog.Any("err", err))
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	pvzID := req.PvzId.String()

	receptionID, err := s.Service.Reception.GetLastReceptionID(ctx, pvzID)
	if errors.Is(err, er.ErrNoOpenReception) {
		slog.Error("no active reception", slog.Any("err", err))
		http.Error(w, `{"message":"no open reception"}`, http.StatusBadRequest)
	}
	if err != nil {
		slog.Error("failed to get reception", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	productID := uuid.New()

	product := models.Product{
		ID:          productID.String(),
		DateTime:    now,
		Type:        string(req.Type),
		ReceptionID: receptionID,
	}

	err = s.Service.Product.AddProduct(r.Context(), product)
	if err != nil {
		slog.Error("failed to add product", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}
	slog.Info("product has been created", slog.Any("info:", product))

	openapiID := openapi_types.UUID(productID)
	resp := openapi.Product{
		Id:          &openapiID,
		DateTime:    &now,
		Type:        openapi.ProductType(req.Type),
		ReceptionId: req.PvzId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) CloseReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzID := chi.URLParam(r, "pvzId")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	receptionID, err := s.Service.Reception.GetOpenReceptionID(ctx, pvzID)
	if errors.Is(err, er.ErrNoOpenReception) {
		slog.Error("no active reception to close", slog.Any("err", err))
		http.Error(w, `{"message":"no open reception to close"}`, http.StatusBadRequest)
	}
	if err != nil {
		slog.Error("failed to get reception ID ", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	err = s.Service.Reception.CloseReception(ctx, receptionID)
	if err != nil {
		slog.Error("failed to close reception", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("reception has been deteted", slog.Any("info:", receptionID))

	openapiID := openapi_types.UUID(uuid.MustParse(receptionID))
	parsedPvzID := openapi_types.UUID(uuid.MustParse(pvzID))

	now := time.Now().UTC()

	resp := openapi.Reception{
		Id:       &openapiID,
		DateTime: now,
		PvzId:    openapi_types.UUID(parsedPvzID),
		Status:   openapi.Close,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) DeleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzID := chi.URLParam(r, "pvzId")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	err := s.Service.Product.DeleteLastProduct(ctx, pvzID)
	if errors.Is(err, er.ErrNoProducts) {
		http.Error(w, `{"message":"nothing to delete"}`, http.StatusBadRequest)
		return
	}
	if errors.Is(err, er.ErrNoOpenReception) {
		http.Error(w, `{"message":"no open reception for deliting product"}`, http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("failed to delete product", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("last product has been deteted", slog.Any("info:", pvzID))

	w.WriteHeader(http.StatusOK)
}

func (s *Server) ListPVZHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.Cfg.HTTP.Timeout)*time.Millisecond)
	defer cancel()

	var (
		start, end *time.Time
		page       = 1
		limit      = s.Cfg.Limits.PaginationLimit
	)

	if v := q.Get("startDate"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			slog.Error("start date is not parsed", slog.Any("err", err))
			http.Error(w, `{"message":"internal error"}`, http.StatusBadRequest)
			return
		}
		t = t.UTC()
		start = &t
	}

	if v := q.Get("endDate"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			slog.Error("end date is not parsed", slog.Any("err", err))
			http.Error(w, `{"message":"internal error"}`, http.StatusBadRequest)
			return
		}
		t = t.UTC()
		end = &t
	}

	if v := q.Get("page"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil && p <= 0 {
			slog.Error("page is not correct", slog.Any("err", err))
			http.Error(w, `{"message":"internal error"}`, http.StatusBadRequest)
			return
		}
		page = p
	}

	data, err := s.Service.PVZ.ListPVZ(ctx, start, end, page, limit)
	if errors.Is(err, er.ErrNoPVZ) {
		http.Error(w, `{"message":"PVZ not found"}`, http.StatusBadRequest)
		return
	}

	if err != nil {
		slog.Error("failed to list pvz", slog.Any("err", err))
		http.Error(w, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}

	var response []map[string]any

	for _, item := range data {
		pvzID := uuid.MustParse(item.PVZ.ID)
		openapiPVZ := openapi.PVZ{
			Id:               (*openapi_types.UUID)(&pvzID),
			City:             openapi.PVZCity(item.PVZ.City),
			RegistrationDate: &item.PVZ.RegistrationDate,
		}

		var receptions []map[string]any
		for _, rec := range item.Receptions {
			recID := uuid.MustParse(rec.Reception.ID)
			recMap := map[string]any{
				"reception": openapi.Reception{
					Id:       (*openapi_types.UUID)(&recID),
					DateTime: rec.Reception.DateTime,
					PvzId:    openapi_types.UUID(pvzID),
					Status:   openapi.ReceptionStatus(rec.Reception.Status),
				},
			}

			var products []openapi.Product
			for _, p := range rec.Products {
				pID := uuid.MustParse(p.ID)
				rID := uuid.MustParse(p.ReceptionID)
				products = append(products, openapi.Product{
					Id:          (*openapi_types.UUID)(&pID),
					DateTime:    &p.DateTime,
					ReceptionId: openapi_types.UUID(rID),
					Type:        openapi.ProductType(p.Type),
				})
			}

			recMap["products"] = products
			receptions = append(receptions, recMap)
		}

		response = append(response, map[string]any{
			"pvz":        openapiPVZ,
			"receptions": receptions,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(s.LoggingMiddleware)
	router.Use(s.PrometheusMiddleware)

	router.Mount("/swagger", api.Routes(router))
	router.Post("/register", s.RegisterHandler)
	router.Post("/login", s.LoginHandler)
	router.Post("/dummyLogin", s.DummyLoginHandler)

	router.Group(func(protected chi.Router) {
		protected.Use(s.RequireAuth)
		protected.Get("/pvz", s.ListPVZHandler)

		employee := protected.With(RequireRole("employee"))
		employee.Post("/products", s.AddProductHandler)
		employee.Post("/pvz/{pvzId}/close_last_reception", s.CloseReceptionHandler)
		employee.Post("/pvz/{pvzId}/delete_last_product", s.DeleteLastProductHandler)
		employee.Post("/receptions", s.CreateReceptionHandler)

		moderator := protected.With(RequireRole("moderator"))
		moderator.Post("/pvz", s.CreatePVZHandler)
	})

	return router
}
