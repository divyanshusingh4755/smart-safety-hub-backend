package shared

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JwtManager struct {
	publicKey *rsa.PublicKey
	logger    *zap.Logger
}

func NewJWTManager(pubPath string, logger *zap.Logger) (*JwtManager, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPath))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse public key: %v", err)
	}

	return &JwtManager{
		publicKey: pubKey,
		logger:    logger,
	}, nil
}

func (m *JwtManager) Verify(tokenString string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userClaims := &UserClaims{
			UserID: claims["sub"].(string),
			Role:   claims["role"].(string),
		}

		if permStr, ok := claims["permissions"].(string); ok {
			trimmed := strings.Trim(permStr, "{}")
			if trimmed != "" {
				userClaims.Permissions = strings.Split(trimmed, ",")
			} else {
				userClaims.Permissions = []string{}
			}
		}
		return userClaims, nil
	}
	return nil, fmt.Errorf("Invalid token claims")
}
