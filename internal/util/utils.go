package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/RobinHoodArmyHQ/robin-api/pkg/nanoid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type JWTData struct {
	*jwt.RegisteredClaims
	UserId    string   `json:"user_id,omitempty"`
	UserRoles []string `json:"user_roles,omitempty"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwt(userID nanoid.NanoID) (string, error) {
	// set expires at after 3 months
	expiresAt := time.Now().AddDate(0, 3, 0)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTData{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserId:    userID.String(),
		UserRoles: []string{"robin"},
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("auth.jwt_secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJwt(signedToken string) (*JWTData, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("auth.jwt_secret")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTData)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("error parsing claims")
	}

	if time.Now().Unix() > expiresAt.Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func GenerateOtp(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	min := 100000
	otp := min + r.Intn(899999)

	return fmt.Sprintf("%d", otp)
}
