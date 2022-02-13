package identification

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(service Service, logger log.Logger) (*Server, error) {
	var server Server
	server.mux = http.NewServeMux()

	server.mux.Handle("/identify", makeIdentifyHandler(service, logger))

	return &server, nil
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func makeIdentifyHandler(service Service, logger log.Logger) *httptransport.Server {
	var identify endpoint.Endpoint
	identify = makeIdentifyEndpoint(service)
	identify = loggingMiddleware(log.With(logger, "method", "identify"))(identify)
	return httptransport.NewServer(identify, decodeIdentifyRequest, encodeIdentifyResponse)
}
