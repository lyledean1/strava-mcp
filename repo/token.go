package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"stravamcp/model"
	"stravamcp/pkg/client"
)

type TokenRepo interface {
	Get() (*model.RedirectTokenResponse, error)
}
type tokenRepo struct {
	stravaClient     client.StravaClient
	token            *model.RedirectTokenResponse
	clientID         string
	clientSecret     string
	folderPath       string
	refreshTokenFile string
}

func NewTokenRepo(stravaClient client.StravaClient, clientID, clientSecret, folderPath, refreshTokenFile string) TokenRepo {
	return &tokenRepo{stravaClient: stravaClient, clientID: clientID, clientSecret: clientSecret, folderPath: folderPath, refreshTokenFile: refreshTokenFile}
}

func (t *tokenRepo) Get() (*model.RedirectTokenResponse, error) {
	if t.token == nil {
		savedToken, err := t.load()
		if err != nil {
			return nil, err
		}
		t.token = savedToken
	}
	if t.token.IsExpired() {
		newToken, err := t.stravaClient.RefreshToken(t.clientID, t.clientSecret, t.token.RefreshToken)
		if err != nil {
			return nil, err
		}
		t.token.AccessToken = newToken.AccessToken
		t.token.RefreshToken = newToken.RefreshToken
		t.token.ExpiresAt = newToken.ExpiresAt
		t.token.ExpiresIn = newToken.ExpiresIn
		err = Save(t.token, t.refreshTokenFile)
		if err != nil {
			return nil, err
		}
	}
	return t.token, nil
}

func (t *tokenRepo) load() (*model.RedirectTokenResponse, error) {
	var redirectTokenResponse *model.RedirectTokenResponse
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", t.folderPath, t.refreshTokenFile))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &redirectTokenResponse)
	if err != nil {
		return nil, err
	}
	return redirectTokenResponse, err
}

func Save(redirectTokenResponse *model.RedirectTokenResponse, refreshTokenFile string) error {
	jsonData, err := json.MarshalIndent(redirectTokenResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %w", err)
	}
	err = os.WriteFile(refreshTokenFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
