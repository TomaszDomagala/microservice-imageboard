package thread

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type ServiceEndpoints struct {
	PostCommentEndpoint  endpoint.Endpoint
	GetCommentEndpoint   endpoint.Endpoint
	DeleteThreadEndpoint endpoint.Endpoint
	CreateThreadEndpoint endpoint.Endpoint
}

func MakeServiceEndpoints(s Service) ServiceEndpoints {
	return ServiceEndpoints{
		PostCommentEndpoint:  MakePostCommentEndpoint(s),
		GetCommentEndpoint:   MakeGetCommentEndpoint(s),
		DeleteThreadEndpoint: MakeDeleteThreadEndpoint(s),
		CreateThreadEndpoint: MakeCreateThreadEndpoint(s),
	}
}

func MakeDeleteThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteThreadRequest)
		err := s.DeleteThread(req.id)
		return basicErrorResponse{Err: err}, nil
	}
}

func MakeCreateThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createThreadRequest)
		id, err := s.CreateThread(req.Ip, req.Board, req.Body)
		if err != nil {
			return basicErrorResponse{Err: err}, nil
		}
		return createThreadResponse{Id: id}, nil
	}
}

func MakeGetCommentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCommentRequest)
		comment, err := s.GetComment(req.ThreadID, req.CommentID)
		return getCommentResponse{Comment: comment, Error: err}, nil
	}
}

func MakePostCommentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postCommentRequest)
		newId, err := s.PostComment(req.Ip, req.ThreadID, req.Body, req.ParentId)
		return postCommentResponse{Id: int(newId), Error: err}, nil
	}
}

type createThreadRequest struct {
	Ip    string `json:"ip"`
	Board string `json:"board"`
	Body  string `json:"body"`
}
type deleteThreadRequest struct {
	id ThreadID
}

type createThreadResponse struct {
	Id int `json:"id"`
}

type basicErrorResponse struct {
	Err error
}

func (r basicErrorResponse) error() error { return r.Err }

type getCommentRequest struct {
	ThreadID  int `json:"threadID"`
	CommentID int `json:"Id"`
}

type getCommentResponse struct {
	Comment Comment `json:"comment"`
	Error   error   `json:"error,omitempty"`
}

type postCommentRequest struct {
	Ip       string `json:"ip"`
	ThreadID int    `json:"threadID"`
	Body     string `json:"body,omitempty"`
	ParentId int    `json:"parentId"`
}

type postCommentResponse struct {
	Id    int   `json:"Id"`
	Error error `json:"error,omitempty"`
}
