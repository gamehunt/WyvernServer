package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"wyvern/server/internal/transport/http/handler"

	"github.com/bytemare/opaque"
)

func register(identity string, password string) error {
	registerEndpoint, err := url.JoinPath(serverUrl, "register")
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
	credId      := responseData.CredId

	fmt.Printf("Received serverIdent=%s and credId=%s\n", serverIdent, credId)

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
		Identity: &clientIdent,
		CredId:   &credId,
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

func login(identity string, password string) error {
	loginEndpoint, err := url.JoinPath(serverUrl, "login")
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
		SessionId: responseData.SessionId,
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

	clientSessionKey := client.SessionKey()

	fmt.Printf("Logged in with sessionId=%x\n", clientSessionKey)

	return nil
}
