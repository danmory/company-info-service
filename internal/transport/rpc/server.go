package rpc

import (
	"context"
	"errors"
	"log"

	"github.com/danmory/company-info-service/internal/service"
)

type companySearchServer struct {
	UnimplementedCompanyInfoSearcherServer
}

func (s *companySearchServer) Search(ctx context.Context, r *SearchRequest) (*SearchResponse, error) {
	isValidINN, err := service.IsINNValid(r.GetInn())
	if err != nil {
		log.Println("failed to check validity of " + r.GetInn() + ", error: " + err.Error())
		return nil, err
	}
	if !isValidINN {
		log.Println("failed to search for " + r.GetInn() + ", error: invalid inn")
		return nil, errors.New("invalid inn")
	}
	resp, err := service.SearchCompanyByINN(r.GetInn())
	if err != nil {
		log.Println("failed to search for " + r.GetInn() + "error: " + err.Error())
		return nil, err
	}
	if len(resp.Ul) == 0 {
		return nil, errors.New("unknown inn")
	}
	var parsedINN string
	if parsedINN, err = service.ParseRecievedINN(resp.Ul[0].Inn); err != nil {
		return nil, err
	}
	if parsedINN != r.GetInn() {
		return nil, errors.New("cannot find inn")
	}
	log.Println("search for " + r.GetInn() + " succeeded")
	company := resp.Ul[0]
	return &SearchResponse{Inn: parsedINN, Ceo: company.CeoName, Name: company.Name}, nil
}

func NewServer() *companySearchServer {
	return new(companySearchServer)
}
