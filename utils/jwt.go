package utils

import (
	"backend-file-management/constant"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Create_token(userId uint, username string, role string) (string, error) {
	// create the claims
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().AddDate(0, 0, 7).Unix()

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(constant.SECRET_JWT))

	if err != nil {
		return "broken_token", err
	}

	return tokenString, nil
}

func Get_role_from_token(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		// Return the secret key used for signing
		return []byte(constant.SECRET_JWT), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the role claim
		if role, ok := claims["role"].(string); ok {
			return role, nil
		}
	}

	return "", errors.New("unable to retrieve data from token")
}

func Get_username_from_token(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		// Return the secret key used for signing
		return []byte(constant.SECRET_JWT), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the username claim
		if username, ok := claims["username"].(string); ok {
			return username, nil
		}
	}

	return "", errors.New("unable to retrieve data from token")
}

func Get_user_id_from_token(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		// Return the secret key used for signing
		return []byte(constant.SECRET_JWT), nil
	})

	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user_id claim
		if user_id, ok := claims["user_id"].(float64); ok {
			return uint(user_id), nil
		}
	}

	return 0, errors.New("unable to retrieve user_id from token")
}
