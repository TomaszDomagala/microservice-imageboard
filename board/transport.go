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
	boardID BoardID `json:"id"`
}

type getThreadResponse struct {
	ids []ThreadID `json:"ids"`
	Err error      `json:"err,omitempty"`
}

func (r getThreadResponse) error() error {
	return r.Err
}

func MakeGetThreadsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getThreadsRequest)
		ids, err := s.GetThreads(req.boardID)
		return getThreadResponse{ids: ids, Err: err}, nil
	}
}

type deleteThreadRequest struct {
	board  BoardID  `json:"boardID"`
	thread ThreadID `json:"threadID"`
}

type basicErrorResponse struct {
	Err error
}

func (r basicErrorResponse) error() error { return r.Err }

func MakeDeleteThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteThreadRequest)
		err := s.DeleteThread(req.board, req.thread)
		return basicErrorResponse{Err: err}, nil
	}
}

type createThreadRequest struct {
	board BoardID `json:"boardID"`
	owner UserID  `json:"owner"`
}

type createThreadResponse struct {
	threadID ThreadID `json:"id"`
	Err      error    `json:"err,omitempty"`
}

func (r createThreadResponse) error() error {
	return r.Err
}

func MakeCreateThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createThreadRequest)
		id, err := s.CreateThread(req.board, req.owner)
		return createThreadResponse{id, err}, nil
	}
}
