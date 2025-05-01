package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Global Variables untuk Expiry Time
const JWT_EXP_MIN = 30          // Token akses expired dalam 5 menit
const REF_EXP_MIN = 7 * 24 * 60 // Token refresh expired dalam 7 hari

// Secret key untuk signing dan verifying JWT
type JwtService struct {
	jwtSecretKey     []byte
	refreshSecretKey []byte
}

// Claims structure untuk payload JWT
type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtService JwtService

// buat mutex untuk mengunci akses ke jwtService

func JwtInit(jwtSecret string, refreshSecret string) {

	jwtService.jwtSecretKey = []byte(jwtSecret)
	jwtService.refreshSecretKey = []byte(refreshSecret)
}

func JwtPrint() {

	fmt.Printf("JWT Secret %s\n", jwtService.jwtSecretKey)
	fmt.Printf("Refresh Secret %s\n", jwtService.refreshSecretKey)

}

// GenerateAccessToken membuat token akses (JWT biasa)
func GenerateAccessToken(user_id string) (string, error) {

	claims := &Claims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(JWT_EXP_MIN))),
			Issuer:    "aethergrow",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	generatedToken, err := token.SignedString(jwtService.jwtSecretKey)

	return generatedToken, err
}

// GenerateRefreshToken membuat token refresh
func GenerateRefreshToken(user_id string) (string, error) {

	claims := &Claims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(REF_EXP_MIN))),
			Issuer:    "aethergrow",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	generatedToken, err := token.SignedString(jwtService.refreshSecretKey)

	return generatedToken, err
}

// ParseAccessToken untuk memverifikasi dan mengurai JWT akses
func ParseAccessToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtService.jwtSecretKey, nil
	})

	if err != nil {
		return nil, &UnauthorizedError{Message: "invalid refresh token"}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, &UnauthorizedError{Message: "invalid refresh token"}
}

// ParseRefreshToken untuk memverifikasi dan mengurai JWT refresh
func ParseRefreshToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtService.refreshSecretKey, nil
	})
	if err != nil {
		return nil, &UnauthorizedError{Message: "invalid refresh token"}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, &UnauthorizedError{Message: "invalid refresh token"}
}

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func PrintJWTInfo(data Claims) {
	fmt.Printf("Id: %s\n", data.UserId)
	fmt.Printf("Issuer: %s\n", data.Issuer)
	fmt.Printf("ExpiresAt: %s\n", data.ExpiresAt)
}
