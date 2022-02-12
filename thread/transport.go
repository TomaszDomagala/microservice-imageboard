package thread

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type ThreadServiceEndpoints struct {
	PostCommentEndpoint endpoint.Endpoint
	GetCommentEndpoint  endpoint.Endpoint
}

func MakeServiceEndpoints(s Service) ThreadServiceEndpoints {
	return ThreadServiceEndpoints{
		PostCommentEndpoint: MakePostCommentEndpoint(s),
		GetCommentEndpoint:  MakeGetCommentEndpoint(s),
	}
}

func MakeGetCommentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCommentRequest)
		comment, err := s.GetComment(req.Id)
		return getCommentResponse{Comment: comment, Error: err}, nil
	}
}

func MakePostCommentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postCommentRequest)
		newId, err := s.PostComment(req.Body, req.Author, req.ParentId)
		return postCommentResponse{Id: int(newId), Error: err}, nil
	}
}

type getCommentRequest struct {
	Id int `json:"id"`
}

type getCommentResponse struct {
	Comment Comment `json:"comment"`
	Error   error   `json:"error,omitempty"`
}

type postCommentRequest struct {
	Body     string `json:"body,omitempty"`
	Author   string `json:"author,omitempty"`
	ParentId int    `json:"parentId"`
}

type postCommentResponse struct {
	Id    int   `json:"id"`
	Error error `json:"error,omitempty"`
}

