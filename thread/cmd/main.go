package main


import (
	"flag"
	"github.com/TomaszDomagala/microservice-imageboard/thread"
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

	svc := thread.NewInMemoryService()
	svc = thread.ServiceLoggingMiddleware(logger)(svc)

	endpoints := thread.MakeServiceEndpoints(svc)

	postCommentHandler := httptransport.NewServer(
		endpoints.PostCommentEndpoint,
		thread.DecodePostCommentRequest,
		thread.EncodeResponse,
	)
	getCommentHandler := httptransport.NewServer(
		endpoints.GetCommentEndpoint,
		thread.DecodeGetCommentRequest,
		thread.EncodeResponse,
	)

	http.Handle("/postComment", postCommentHandler)
	http.Handle("/getComment", getCommentHandler)

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}

