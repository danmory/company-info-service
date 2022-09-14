package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/danmory/company-info-service/internal/core"
)

var (
	// use Sprintf to add INN into search url
	sourceAPI = "https://www.rusprofile.ru/ajax.php?query=%v&action=search" 
	innRegexp = `^(\d{10}|\d{12})$`
)

func constructSearchAddress(inn string) string {
	return fmt.Sprintf(sourceAPI, inn)
}

func IsINNValid(inn string) (bool, error) {
	return regexp.MatchString(innRegexp, inn)
}

func SearchCompanyByINN(inn string) (*core.SearchAPIResponse, error) {
	resp, err := http.Get(constructSearchAddress(inn))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result core.SearchAPIResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
