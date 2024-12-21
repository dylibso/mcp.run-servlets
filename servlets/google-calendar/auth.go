// auth.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	pdk "github.com/extism/go-pdk"
)

const (
	deviceCodeEndpoint = "https://oauth2.googleapis.com/device/code"
	tokenEndpoint      = "https://oauth2.googleapis.com/token"
)

type AuthManager struct {
	clientID     string
	clientSecret string
}

type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error,omitempty"`
}

func NewAuthManager() (*AuthManager, error) {
	clientID, ok := pdk.GetConfig("client_id")
	if !ok {
		return nil, errors.New("client_id config required")
	}

	clientSecret, ok := pdk.GetConfig("client_secret")
	if !ok {
		return nil, errors.New("client_secret config required")
	}

	return &AuthManager{
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

func (am *AuthManager) StartDeviceFlow() (*DeviceCode, error) {
	data := map[string]string{
		"client_id": am.clientID,
		"scope":     "https://www.googleapis.com/auth/calendar",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req := pdk.NewHTTPRequest(pdk.MethodPost, deviceCodeEndpoint)
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(body)

	res := req.Send()
	if res.Status() != 200 {
		return nil, fmt.Errorf("device code request failed: %s", string(res.Body()))
	}

	var deviceCode DeviceCode
	if err := json.Unmarshal(res.Body(), &deviceCode); err != nil {
		return nil, fmt.Errorf("failed to parse device code response: %v", err)
	}

	return &deviceCode, nil
}

func (am *AuthManager) PollForToken(deviceCode string) (*TokenResponse, error) {
	req := pdk.NewHTTPRequest(pdk.MethodPost, tokenEndpoint)
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	// Use form-urlencoded format for token request
	body := fmt.Sprintf(
		"client_id=%s&client_secret=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code",
		am.clientID, am.clientSecret, deviceCode,
	)
	req.SetBody([]byte(body))

	res := req.Send()

	var tokenResp TokenResponse
	if err := json.Unmarshal(res.Body(), &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	// If authorization pending or slow down, return the error
	if tokenResp.Error == "authorization_pending" || tokenResp.Error == "slow_down" {
		return &tokenResp, nil
	}

	// For other errors, or if status not 200, return error
	if res.Status() != 200 || tokenResp.Error != "" {
		return nil, fmt.Errorf("token request failed: %s", tokenResp.Error)
	}

	return &tokenResp, nil
}

func (am *AuthManager) RefreshToken(refreshToken string) (*TokenResponse, error) {
	data := map[string]string{
		"client_id":     am.clientID,
		"client_secret": am.clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req := pdk.NewHTTPRequest(pdk.MethodPost, tokenEndpoint)
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(body)

	res := req.Send()

	if res.Status() != 200 {
		return nil, fmt.Errorf("token refresh failed")
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(res.Body(), &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &tokenResp, nil
}
