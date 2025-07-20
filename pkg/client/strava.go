package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stravamcp/model"
	"strings"
)

const oauthTokenUrl = "%s/oauth/token"

type StravaClient interface {
	GetTokenFromAuthCode(clientID, clientSecret, authorizationCode string) (*model.RedirectTokenResponse, error)
	RefreshToken(clientID, clientSecret, refreshToken string) (*model.TokenResponse, error)
	GetAthleteActivityByID(id, accessToken string) (*model.AthleteActivity, error)
	GetAllAthleteActivities(after int, accessToken string) ([]model.AthleteActivity, error)
	FetchStreams(activityID string, keys []string, accessToken string) (*model.ActivityStreams, error)
}

type stravaClient struct {
	client       *http.Client
	baseUrl      string
	perPageLimit int // maximum number of activities per page
}

func NewStravaClient(baseUrl string) StravaClient {
	return &stravaClient{client: &http.Client{}, baseUrl: baseUrl, perPageLimit: 200}
}

func (s *stravaClient) makeRequest(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

func (s *stravaClient) makeJSONRequest(method, url string, body io.Reader, headers map[string]string, target interface{}) error {
	resp, err := s.makeRequest(method, url, body, headers)
	if err != nil {
		return err
	}

	//nolint: errcheck // defer handles after function exit
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}

func (s *stravaClient) makeAuthenticatedRequest(method, url, accessToken string, target interface{}) error {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	return s.makeJSONRequest(method, url, nil, headers, target)
}

func (s *stravaClient) GetTokenFromAuthCode(clientID, clientSecret, authorizationCode string) (*model.RedirectTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", authorizationCode)
	data.Set("grant_type", "authorization_code")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	targetUrl := fmt.Sprintf(oauthTokenUrl, s.baseUrl)
	var tokenResponse model.RedirectTokenResponse

	err := s.makeJSONRequest("POST", targetUrl, strings.NewReader(data.Encode()), headers, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (s *stravaClient) RefreshToken(clientID, clientSecret, refreshToken string) (*model.TokenResponse, error) {
	tokenRequest := model.TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}

	jsonBody, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	var tokenResponse model.TokenResponse
	err = s.makeJSONRequest("POST", "https://www.strava.com/oauth/token", strings.NewReader(string(jsonBody)), headers, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (s *stravaClient) GetAthleteActivityByID(id, accessToken string) (*model.AthleteActivity, error) {
	athleteActivityUrl := fmt.Sprintf("%s/api/v3/athlete/activities/%s", s.baseUrl, id)

	var activity model.AthleteActivity
	err := s.makeAuthenticatedRequest("GET", athleteActivityUrl, accessToken, &activity)
	if err != nil {
		return nil, fmt.Errorf("fetching activities: %w", err)
	}

	return &activity, nil
}

func (s *stravaClient) GetAthleteActivity(after, page int, accessToken string) ([]model.AthleteActivity, error) {
	athleteActivityUrl := fmt.Sprintf("%s/api/v3/athlete/activities?page=%d&per_page=%d&after=%d",
		s.baseUrl, page, s.perPageLimit, after)

	var activities []model.AthleteActivity
	err := s.makeAuthenticatedRequest("GET", athleteActivityUrl, accessToken, &activities)
	if err != nil {
		return nil, fmt.Errorf("fetching activities: %w", err)
	}

	return activities, nil
}

func (s *stravaClient) GetAllAthleteActivities(after int, accessToken string) ([]model.AthleteActivity, error) {
	var allActivities []model.AthleteActivity
	page := 1

	for {
		activities, err := s.GetAthleteActivity(after, page, accessToken)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", page, err)
		}

		if len(activities) == 0 {
			break
		}
		allActivities = append(allActivities, activities...)
		page++
	}

	return allActivities, nil
}

func (s *stravaClient) FetchStreams(activityID string, keys []string, accessToken string) (*model.ActivityStreams, error) {
	keysStr := strings.Join(keys, ",")
	url := fmt.Sprintf(
		"https://www.strava.com/api/v3/activities/%s/streams?keys=%s&key_by_type=true",
		activityID,
		keysStr,
	)

	var streams model.ActivityStreams
	err := s.makeAuthenticatedRequest("GET", url, accessToken, &streams)
	if err != nil {
		return nil, fmt.Errorf("fetching streams: %w", err)
	}

	return &streams, nil
}
