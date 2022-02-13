package media

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"io"
)

type ServiceEndpoints struct {
	PostMediaEndpoint   endpoint.Endpoint
	GetMediaEndpoint    endpoint.Endpoint
	DeleteMediaEndpoint endpoint.Endpoint
}

func MakeServiceEndpoints(s Service) ServiceEndpoints {
	return ServiceEndpoints{
		PostMediaEndpoint:   MakePostMediaEndpoint(s),
		GetMediaEndpoint:    MakeGetMediaEndpoint(s),
		DeleteMediaEndpoint: MakeDeleteMediaEndpoint(s),
	}
}

func MakePostMediaEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postMediaRequest)
		err := s.PostMedia(req.Name, req.Reader)
		return postMediaResponse{err}, nil
	}
}

func MakeGetMediaEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getMediaRequest)
		data, err := s.GetMedia(req.Name)
		return getMediaResponse{data, err}, nil
	}
}

func MakeDeleteMediaEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteMediaRequest)
		e := s.DeleteMedia(req.Name)
		return deleteMediaResponse{e}, nil
	}
}

type deleteMediaRequest struct {
	Name string `json:"name"`
}

type deleteMediaResponse struct {
	Error error `json:"error,omitempty"`
}

func (r deleteMediaResponse) error() error { return r.Error }

type getMediaRequest struct {
	Name string `json:"name"`
}

type getMediaResponse struct {
	Data  []byte
	Error error `json:"error,omitempty"`
}

func (r getMediaResponse) error() error { return r.Error }

type postMediaRequest struct {
	Name   string
	Reader io.Reader
}

type postMediaResponse struct {
	Error error `json:"error,omitempty"`
}

func (r postMediaResponse) error() error { return r.Error }
