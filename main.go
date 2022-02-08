package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	svc := NewInMemoryService()
	svc = threadServiceloggingMiddleware(logger)(svc)

	endpoints := makeThreadServiceEndpoints(svc)

	postCommentHandler := httptransport.NewServer(
		endpoints.PostCommentEndpoint,
		decodePostCommentRequest,
		encodeResponse,
	)
	getCommentHandler := httptransport.NewServer(
		endpoints.GetCommentEndpoint,
		decodeGetCommentRequest,
		encodeResponse,
	)

	http.Handle("/postComment", postCommentHandler)
	http.Handle("/getComment", getCommentHandler)

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
