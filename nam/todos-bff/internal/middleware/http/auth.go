package http_middleware

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey             contextKey = "user_id"
	accessTokenCookieName string     = "access_token"
)

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (int64, error)
}

type JWTTokenVerifier struct {
	secret []byte
}

func NewJWTTokenVerifier() (*JWTTokenVerifier, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	return &JWTTokenVerifier{secret: []byte(secret)}, nil
}

func (v *JWTTokenVerifier) VerifyAccessToken(_ context.Context, accessToken string) (int64, error) {
	if v == nil {
		return 0, errors.New("token verifier is nil")
	}
	if strings.TrimSpace(accessToken) == "" {
		return 0, errors.New("access token is required")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				alg := "unknown"
				if token.Method != nil {
					alg = token.Method.Alg()
				}
				return nil, fmt.Errorf("unexpected signing method: %s", alg)
			}
			return v.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return 0, fmt.Errorf("parse access token: %w", err)
	}
	if !token.Valid {
		return 0, errors.New("invalid access token")
	}

	userID, err := parseUserIDFromClaims(claims)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func AuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler {
	if verifier == nil {
		panic("auth middleware requires token verifier")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(accessTokenCookieName)
			if err != nil || strings.TrimSpace(cookie.Value) == "" {
				next.ServeHTTP(w, r)
				return
			}

			userID, err := verifier.VerifyAccessToken(r.Context(), cookie.Value)
			if err != nil {
				ClearAccessTokenCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			ctx := WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithUserID(ctx context.Context, userID int64) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if userID <= 0 {
		return ctx
	}

	return context.WithValue(ctx, UserIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	if ctx == nil {
		return 0, false
	}

	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok || userID <= 0 {
		return 0, false
	}

	return userID, true
}

func ClearAccessTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func parseUserIDFromClaims(claims jwt.MapClaims) (int64, error) {
	if raw, ok := claims["user_id"]; ok {
		userID, err := parseUserIDClaim(raw)
		if err == nil {
			return userID, nil
		}
	}

	if raw, ok := claims["sub"]; ok {
		userID, err := parseUserIDClaim(raw)
		if err == nil {
			return userID, nil
		}
	}

	return 0, errors.New("user id claim is missing or invalid")
}

func parseUserIDClaim(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		if v <= 0 {
			return 0, errors.New("user id must be positive")
		}
		return v, nil
	case int:
		if v <= 0 {
			return 0, errors.New("user id must be positive")
		}
		return int64(v), nil
	case float64:
		if v <= 0 || v != math.Trunc(v) {
			return 0, errors.New("user id must be a positive integer")
		}
		return int64(v), nil
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user id claim: %w", err)
		}
		if parsed <= 0 {
			return 0, errors.New("user id must be positive")
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported user id claim type: %T", value)
	}
}
