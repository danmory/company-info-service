package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/danmory/company-info-service/internal/core"
)

var (
	// use Sprintf to add INN into search url
	sourceAPI         = "https://www.rusprofile.ru/ajax.php?query=%v&action=search"
	innRegexp         = `^(\d{10}|\d{12})$`
	innRetrievingRegexp = `(\d{10}|\d{12})`
	innRetriever *regexp.Regexp
)

func constructSearchAddress(inn string) string {
	return fmt.Sprintf(sourceAPI, inn)
}

// inn returns from API with extra symbols like ~ and !
func ParseRecievedINN(received string) (string, error) {
	if innRetriever == nil {
		if r, err := regexp.Compile(innRetrievingRegexp); err != nil {
			log.Println("error compiling inn regexp: " + err.Error())
			return "", errors.New("internal error")
		} else {
			innRetriever = r
		}
	}
	return string(innRetriever.Find([]byte(received))), nil
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
