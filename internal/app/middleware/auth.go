package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ruslanjo/url_shortener/internal/config"
)

var ErrTokenNotValid = errors.New("invalid JWT token")

type CustomClaims struct {
	UserID string
}

type Claims struct {
	CustomClaims
	jwt.RegisteredClaims
}

func Signup(next http.HandlerFunc, tokenGen TokenGenerator) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var userID string

		cookie, err := r.Cookie(config.AuthCookie)

		switch {
		case err == nil:
			tokenString := cookie.Value
			claims, err := tokenGen.GetClaims(tokenString)

			if err != nil {
				err = AddAuthCookie(&userID, tokenGen, w, r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			} else {
				userID = claims.UserID
			}

		case errors.Is(err, http.ErrNoCookie):
			err = AddAuthCookie(&userID, tokenGen, w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), config.CtxUserIDKey, userID)
		next(w, r.WithContext(ctx))
	}

	return fn
}

type TokenGenerator struct{}

func (t TokenGenerator) Create(customClaims CustomClaims) (string, error) {
	claims := Claims{
		CustomClaims: customClaims,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	JWTString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}
	return JWTString, nil
}

func (t TokenGenerator) GetClaims(tokenString string) (Claims, error) {
	claims := Claims{}

	secretFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenNotValid
		}
		return []byte(config.JWTSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, secretFunc)
	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, ErrTokenNotValid
	}

	return claims, nil
}

func generateUUID() string {
	id := uuid.New()
	return id.String()
}

func AddAuthCookie(
	userID *string,
	tokenGen TokenGenerator,
	w http.ResponseWriter,
	r *http.Request,
) error {
	*userID = generateUUID()

	claims := CustomClaims{
		UserID: *userID,
	}

	tokenString, err := tokenGen.Create(claims)
	if err != nil {
		return err
	}

	newCookie := http.Cookie{
		Name:  config.AuthCookie,
		Value: tokenString,
	}
	http.SetCookie(w, &newCookie)

	return nil
}
