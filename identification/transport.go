package identification

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type identifyRequest struct {
	Ip string `json:"ip"`
}

type identifyResponse struct {
	Id string `json:"id"`
}

func makeIdentifyEndpoint(srv Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identifyRequest)
		id, err := srv.Identify(req.Ip)
		if err != nil {
			return nil, err
		}
		return identifyResponse{id}, nil
	}
}


