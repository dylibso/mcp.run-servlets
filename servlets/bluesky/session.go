package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/extism/go-pdk"
)

type LoginPayload struct {
	Handle   string `json:"identifier"`
	Password string `json:"password"`
}

type Session struct {
	AccessJwt string `json:"accessJwt"`
	DID       string `json:"did"`
}

func loadConfig() error {
	BASE_URL, _ = pdk.GetConfig("BASE_URL") // default https://bsky.social
	HANDLE, _ = pdk.GetConfig("HANDLE")
	PASSWORD, _ = pdk.GetConfig("APP_PASSWORD")

	if BASE_URL == "" {
		BASE_URL = "https://bsky.social"
	}
	if HANDLE == "" {
		return errors.New("handle is required")
	}
	if PASSWORD == "" {
		return errors.New("password is required")
	}
	return nil
}

func loginSession() error {
	url := BASE_URL + "/xrpc/com.atproto.server.createSession"
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Content-Type", "application/json")
	loginPayload := LoginPayload{
		Handle:   HANDLE,
		Password: PASSWORD,
	}
	jsonBytes, err := json.Marshal(&loginPayload)
	if err != nil {
		return err
	}
	req.SetBody(jsonBytes)
	resp := req.Send()
	if resp.Status() != http.StatusOK {
		return fmt.Errorf("failed to login: %d, %s", resp.Status(), string(resp.Body()))
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &currentSession); err != nil {
		return err
	}
	pdk.Log(pdk.LogInfo, "logged in")
	return nil
}
