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
		if t.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.publicKey, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("missing or invalid sub claim")
	}

	// Role may be absent in password-reset tokens.
	role, _ := claims["role"].(string)

	userClaims := &UserClaims{
		UserID:      userID,
		Role:        role,
		Token:       token,
		Permissions: []string{},
	}

	// Permissions may be absent when rolesPermissions is nil.
	switch permissions := claims["permissions"].(type) {
	case string:
		trimmed := strings.Trim(permissions, "{}")
		if trimmed != "" {
			userClaims.Permissions = strings.Split(trimmed, ",")
		}
	case []interface{}:
		for _, permission := range permissions {
			if p, ok := permission.(string); ok {
				userClaims.Permissions = append(
					userClaims.Permissions, p,
				)
			}
		}
	}
	return userClaims, nil
}
