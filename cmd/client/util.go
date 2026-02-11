package main

import (
	"fmt"
	"io"
	"net/http"
)

func readBodyAsString(resp *http.Response) (string, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func checkResponseStatusExt(resp *http.Response, desired int) error {
	if resp.StatusCode != desired {
		message, err := readBodyAsString(resp)
		if err != nil {
			return fmt.Errorf("Unknown failure with status %d", resp.StatusCode)
		} else {
			return fmt.Errorf("Request failed with status %d, message: %s", resp.StatusCode, message)
		}
	}	
	return nil
}

func checkResponseStatus(resp *http.Response) error {
	return checkResponseStatusExt(resp, http.StatusOK)	
}
