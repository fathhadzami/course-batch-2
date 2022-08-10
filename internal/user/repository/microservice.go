package repository

import (
	"context"
	"course/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	getUserURL = "internal/users/"
)

type McsrvRepo struct {
	hostName string
	username string
	password string
	client   *http.Client
}

func NewMcsrvRepo() *McsrvRepo {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	return &McsrvRepo{
		hostName: "http://localhost:8083/",
		username: "user",
		password: "abcd1234",
		client:   &client,
	}
}

func (dr McsrvRepo) IsUserExist(ctx context.Context, userID int) bool {
	url := fmt.Sprintf("%s%s%d", dr.hostName, getUserURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	req.SetBasicAuth(dr.username, dr.password)
	resp, err := dr.client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	var user domain.User
	json.NewDecoder(resp.Body).Decode(&user)
	return user.ID > 0
}
