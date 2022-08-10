package user

import (
	"context"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	privateKey []byte = []byte("mySignaturePrivateKey")
)

type UserRepo interface {
	IsUserExist(ctx context.Context, userID int) bool
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uu UserUsecase) IsUserExist(ctx context.Context, userID int) bool {
	if userID == 0 {
		return false
	}
	return uu.repo.IsUserExist(ctx, userID)
}

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iss":     "edspert",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu UserUsecase) DecriptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth invalid")
		}
		return privateKey, nil
	})

	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}

	if !parsedToken.Valid {
		return data, errors.New("invalid token")
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}
