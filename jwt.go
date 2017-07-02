package gfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"time"
)

var (
	ErrNoAuthHeader      error = errors.New("No authorization header")
	ErrInvalidAuthHeader error = errors.New("Invalid authorization header")
)

func getValidationKeyGetter(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Unexpected signing method: %v", method)
			}
		}
		return secret, nil
	}
}

func GetTokenData(tokenString string, secret []byte, output interface{}) error {
	token, err := jwt.Parse(tokenString, getValidationKeyGetter(secret))
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"]
		if ok {
			err = json.Unmarshal([]byte(sub.(string)), output)
			return err
		} else {
			return errors.New("Invalid token. Sub was not set. ")
		}
	} else {
		return err
	}
}

func GetToken(secret []byte) (string, error) {
	exp := time.Now().Add(31 * 24 * time.Hour)
	claim := &jwt.StandardClaims{
		ExpiresAt: exp.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        uuid.NewV4().String(),
		Subject:   "{}",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(secret)
}
