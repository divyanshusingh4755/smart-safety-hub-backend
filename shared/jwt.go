package shared

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type RolesPermissions struct {
	Role        string `db:"role"`
	Permissions string `db:"permissions"`
}

type JwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	logger     *zap.Logger
}

func NewJWTManager(privPath, pubPath string, logger *zap.Logger) (*JwtManager, error) {
	pvteKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privPath))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPath))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse public key: %v", err)
	}

	return &JwtManager{
		privateKey: pvteKey,
		publicKey:  pubKey,
		logger:     logger,
	}, nil
}

func (m *JwtManager) GenerateToken(userId string, rolesPermissions *RolesPermissions, ttl time.Duration) (string, error) {
	var claims jwt.MapClaims
	if rolesPermissions != nil {
		claims = jwt.MapClaims{
			"sub":         userId,
			"role":        rolesPermissions.Role,
			"permissions": rolesPermissions.Permissions,
			"iat":         time.Now().Unix(),
			"exp":         time.Now().Add(ttl).Unix(),
		}
	} else {
		claims = jwt.MapClaims{
			"sub": userId,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(ttl).Unix(),
		}

	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
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
			Token:  token,
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
