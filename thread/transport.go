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
	GetChildrenEndpoint  endpoint.Endpoint
}

func MakeServiceEndpoints(s Service) ServiceEndpoints {
	return ServiceEndpoints{
		PostCommentEndpoint:  MakePostCommentEndpoint(s),
		GetCommentEndpoint:   MakeGetCommentEndpoint(s),
		GetChildrenEndpoint:  MakeGetChildrenEndpoint(s),
		DeleteThreadEndpoint: MakeDeleteThreadEndpoint(s),
		CreateThreadEndpoint: MakeCreateThreadEndpoint(s),
	}
}

func MakeDeleteThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteThreadRequest)
		err := s.DeleteThread(req.Id)
		return basicErrorResponse{Err: err}, nil
	}
}

func MakeCreateThreadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createThreadRequest)
		id, err := s.CreateThread(req.Ip, req.Board, req.Body, req.HasMedia)
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

func MakeGetChildrenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCommentChildrenRequest)
		comments, err := s.GetChildren(req.ThreadID, req.CommentID)
		return getChildrenResponse{comments, err}, nil
	}
}

func MakePostCommentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postCommentRequest)
		newId, err := s.PostComment(req.Ip, req.ThreadID, req.Body, req.ParentId, req.HasMedia)
		return postCommentResponse{Id: newId, Error: err}, nil
	}
}

type createThreadRequest struct {
	Ip       string `json:"ip"`
	Board    string `json:"board"`
	Body     string `json:"body"`
	HasMedia bool   `json:"hasMedia"`
}

type deleteThreadRequest struct {
	Id ThreadID `json:"id"`
}

type createThreadResponse struct {
	Id int `json:"id"`
}

type basicErrorResponse struct {
	Err error `json:"err,omitempty"`
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

type getCommentChildrenRequest struct {
	ThreadID  int `json:"threadID"`
	CommentID int `json:"Id"`
}

type getChildrenResponse struct {
	Children []Comment `json:"children"`
	Error    error     `json:"error,omitempty"`
}

func (r getChildrenResponse) error() error { return r.Error }

type postCommentRequest struct {
	Ip       string `json:"ip"`
	ThreadID int    `json:"threadID"`
	Body     string `json:"body,omitempty"`
	ParentId int    `json:"parentId"`
	HasMedia bool   `json:"hasMedia,omitempty"`
}

type postCommentResponse struct {
	Id    int   `json:"Id"`
	Error error `json:"error,omitempty"`
}
