package pb

import (
	"context"
	"errors"

	"github.com/danmory/company-info-service/internal/service"
)

type companySearchServer struct {
	UnimplementedCompanyInfoSearcherServer
}

func (s *companySearchServer) Search(ctx context.Context, r *SearchRequest) (*SearchResponse, error) {
	resp, err := service.SearchCompanyByINN(r.GetInn())
	if err != nil {
		return nil, err
	}
	if len(resp.Ul) == 0 {
		return nil, errors.New("unknown inn")
	}
	company := resp.Ul[0]
	return &SearchResponse{Inn: company.Inn, Ceo: company.CeoName, Name: company.Name}, nil
}

func NewServer() *companySearchServer {
	return new(companySearchServer)
}
