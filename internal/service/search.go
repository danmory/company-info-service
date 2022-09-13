package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danmory/company-info-service/internal/core"
)

var sourceAPI string = "https://www.rusprofile.ru/ajax.php?query=%v&action=search"

func constructSearchAddress(inn string) string {
	return fmt.Sprintf(sourceAPI, inn)
}

func SearchCompanyByINN(inn string) (*core.APIResponse, error) {
	resp, err := http.Get(constructSearchAddress(inn))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result core.APIResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
