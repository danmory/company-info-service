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
	log.Println("search for " + r.GetInn() + " succeeded")
	company := resp.Ul[0]
	return &SearchResponse{Inn: r.GetInn(), Ceo: company.CeoName, Name: company.Name}, nil
}

func NewServer() *companySearchServer {
	return new(companySearchServer)
}
