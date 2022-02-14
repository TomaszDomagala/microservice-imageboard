package board

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type ServiceEndpoints struct {
	CreateThreadEndpoint endpoint.Endpoint
	DeleteThreadEndpoint endpoint.Endpoint
	GetThreadsEndpoint   endpoint.Endpoint
}

func MakeServiceEndpoints(s Service) ServiceEndpoints {
	return ServiceEndpoints{
		GetThreadsEndpoint:   MakeGetThreadsEndpoint(s),
		DeleteThreadEndpoint: MakeDeleteThreadEndpoint(s),
		CreateThreadEndpoint: MakeCreateThreadEndpoint(s),
	}
}

type getThreadsRequest struct {
	BoardID BoardID `json:"id"`
}

type getThreadResponse struct {
	Ids []ThreadID `json:"ids"`
	Err error      `json:"err,omitempty"`
}

func (r getThreadResponse) error() error {
	return r.Err
}

func MakeGetThreadsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getThreadsRequest)
		ids, err := s.GetThreads(req.BoardID)
		return getThreadResponse{Ids: ids, Err: err}, nil
	}
}

type deleteThreadRequest struct {
	Board  BoardID  `json:"boardID"`
	Thread ThreadID `json:"threadID"`
}

type basicErrorResponse struct {
	Err error
}

func (r basicErrorResponse) error() error { return r.Err }

func MakeDeleteThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteThreadRequest)
		err := s.DeleteThread(req.Board, req.Thread)
		return basicErrorResponse{Err: err}, nil
	}
}

type createThreadRequest struct {
	Board BoardID `json:"boardID"`
	Owner UserID  `json:"owner"`
}

type createThreadResponse struct {
	ThreadID ThreadID `json:"id"`
	Err      error    `json:"err,omitempty"`
}

func (r createThreadResponse) error() error {
	return r.Err
}

func MakeCreateThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createThreadRequest)
		id, err := s.CreateThread(req.Board, req.Owner)
		return createThreadResponse{id, err}, nil
	}
}
