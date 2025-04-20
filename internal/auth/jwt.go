package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTManager(secret string, expirationTime int) *JWTManager {
	expiration := time.Duration(expirationTime) * time.Minute
	return &JWTManager{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (j *JWTManager) Generate(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString(j.secret)
}

func (j *JWTManager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.Wrap(err, "can't validate token")
	}

	return token.Claims.(*Claims), nil
}
