package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"wyvern/server/internal/transport/http/handler"
	"wyvern/server/internal/types"

	"github.com/bytemare/opaque"
)

func register(identity string, password string) error {
	registerEndpoint, err := url.JoinPath(serverUrl, "auth/register")
	if err != nil {
		return fmt.Errorf("Error creating endpoints: %v", err)
	}

	conf := opaque.DefaultConfiguration()

	client, err := conf.Client()
	if err != nil {
		return fmt.Errorf("Failed to create opaque client: %v", err)
	}

	c1   := client.RegistrationInit([]byte(password))
	msg1 := c1.Serialize()

	reqData := handler.RegisterRequest {
		Payload: msg1,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("Failed to serialize request: %v", err)
	}

	httpclient := &http.Client{}
	resp, err := httpclient.Post(registerEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
	    return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	err = checkResponseStatus(resp)
	if err != nil {
		return err
	}

	var responseData handler.RegisterResponse
	
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return fmt.Errorf("Failed to parse response: %v", err)
	}

	clientIdent := identity
	serverIdent := responseData.Identity
	id          := responseData.UserId

	fmt.Printf("Received serverIdent=%s and userId=%s\n", serverIdent, id)

	response, err := client.Deserialize.RegistrationResponse(responseData.Payload)
	if err != nil {
		return fmt.Errorf("Failed to create registration response: %v", err)
	}

	record, _ := client.RegistrationFinalize(response, opaque.ClientRegistrationFinalizeOptions{
		ClientIdentity: []byte(clientIdent),
		ServerIdentity: serverIdent,
	})
	msg3 := record.Serialize()

	reqData = handler.RegisterRequest {
		UserId:   &id,
		Identity: &clientIdent,
		Payload:   msg3,
	}

	jsonData, err = json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("Failed to serialize request: %v", err)
	}

	resp, err = httpclient.Post(registerEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
	    return fmt.Errorf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	err = checkResponseStatus(resp)
	if err != nil {
		return err
	}

	fmt.Println("Registration successful with clientId=", string(clientIdent))

	return nil
}

var currentSessionKey   []byte
var currentAccessToken  string
var currentRefreshToken string
var currentSessionId    types.ID

func login(identity string, password string) error {
	loginEndpoint, err := url.JoinPath(serverUrl, "auth/login")
	if err != nil {
		return fmt.Errorf("Error creating endpoints: %v", err)
	}

	conf := opaque.DefaultConfiguration()

	client, err := conf.Client()
	if err != nil {
		return fmt.Errorf("Failed to create opaque client: %v", err)
	}

	ke1  := client.LoginInit([]byte(password))
	msg1 := ke1.Serialize()

	reqData := handler.LoginRequest {
		Identity: &identity,
		Payload:   msg1,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("Failed to serialize request: %v", err)
	}


	httpclient := &http.Client{}
	resp, err := httpclient.Post(loginEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
	    return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	err = checkResponseStatus(resp)
	if err != nil {
		return err
	}

	var responseData handler.LoginResponse
	
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return fmt.Errorf("Failed to parse response: %v", err)
	}

	fmt.Printf("Received serverId=%s\n", responseData.Identity)

	ke2, err := client.Deserialize.KE2(responseData.Payload)
	if err != nil {
		return err
	}

	ke3, _, err := client.LoginFinish(ke2, opaque.ClientLoginFinishOptions{
		ClientIdentity: []byte(identity),
		ServerIdentity: responseData.Identity,
	})
	if err != nil {
		return err
	}

	msg3 := ke3.Serialize()

	reqData = handler.LoginRequest {
		OpaqueId: responseData.OpaqueId,
		Payload: msg3,
	}

	jsonData, err = json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("Failed to serialize request: %v", err)
	}

	resp, err = httpclient.Post(loginEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
	    return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	err = checkResponseStatus(resp)
	if err != nil {
		return err
	}

	var finalResp handler.SuccessLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&finalResp)
	if err != nil {
		return fmt.Errorf("Failed to parse response: %v", err)
	}

	currentSessionKey   = client.SessionKey()
	currentSessionId    = finalResp.SessionId
	currentAccessToken  = finalResp.AccessToken
	currentRefreshToken = finalResp.RefreshToken

	fmt.Printf("Logged in with sessionId=%s, accessToken=%s, refreshToken=%s\n", finalResp.SessionId,
	finalResp.AccessToken, finalResp.RefreshToken)

	return nil
}

func refresh() error {
	if currentSessionKey == nil {
		return fmt.Errorf("Not logged in")
	}

	refreshEndpoint, err := url.JoinPath(serverUrl, "auth/refresh")
	if err != nil {
		return fmt.Errorf("Error creating endpoints: %v", err)
	}

	h := sha256.New()
    h.Write(currentSessionKey)
    h.Write([]byte("auth"))
	authKey := h.Sum(nil)

	nonce := opaque.RandomBytes(16)

	h.Reset()
    h.Write(authKey)
    h.Write(nonce)
	proof := h.Sum(nil)

	req := &handler.RefreshRequest{
		SessionId: currentSessionId,
		RefreshToken: currentRefreshToken,
		Proof: proof,	
		Nonce: nonce,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Failed to serialize request: %v", err)
	}

	resp, err := http.Post(refreshEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
	    return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	err = checkResponseStatus(resp)
	if err != nil {
		return err
	}

	var finalResp handler.RefreshResponse
	err = json.NewDecoder(resp.Body).Decode(&finalResp)
	if err != nil {
		return fmt.Errorf("Failed to parse response: %v", err)
	}

	currentAccessToken = finalResp.AccessToken

	fmt.Printf("accessToken=%s\n", currentAccessToken)

	return nil
}
