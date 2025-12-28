package lib

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId   string `json:"user_id"`
	UserName string `json:"username"`
	RoleId   string `json:"role_id"`
	Iss      string `json:"iss"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

func CreateJWT(claim JWTClaims, secretKey string, exp time.Duration, typ string) (string, error) {

	var claims = jwt.MapClaims{
		"iss":       claim.Iss,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(exp).Unix(),
		"typ":       typ,
		"user_id":   claim.UserId,
		"user_name": claim.UserName,
		"role_id":   claim.RoleId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GetUserIdFromAccessToken(token, secretKey string, issuer string) (int64, error) {
	parsedToken, err := jwt.Parse(token,
		func(_ *jwt.Token) (interface{}, error) { return []byte(secretKey), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userIdFloat, ok := claims["user_id"].(string)
		if !ok {
			return 0, fmt.Errorf("user_id not found in token claims")
		}
		if claims["typ"] != "access" {
			return 0, fmt.Errorf("invalid token type")
		}
		var userId int64
		_, err := fmt.Sscan(userIdFloat, &userId)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id format in token claims")
		}
		return userId, nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}

func GetUserIdFromRefreshToken(token, secretKey string, issuer string) (int64, error) {
	parsedToken, err := jwt.Parse(token,
		func(_ *jwt.Token) (interface{}, error) { return []byte(secretKey), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userIdFloat, ok := claims["user_id"].(string)
		if !ok {
			return 0, fmt.Errorf("user_id not found in token claims")
		}
		if claims["typ"] != "refresh" {
			return 0, fmt.Errorf("invalid token type")
		}
		var userId int64
		_, err := fmt.Sscan(userIdFloat, &userId)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id format in token claims")
		}
		return userId, nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}
