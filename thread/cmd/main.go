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
	svc = thread.threadServiceloggingMiddleware(logger)(svc)

	endpoints := thread.makeThreadServiceEndpoints(svc)

	postCommentHandler := httptransport.NewServer(
		endpoints.PostCommentEndpoint,
		thread.decodePostCommentRequest,
		thread.encodeResponse,
	)
	getCommentHandler := httptransport.NewServer(
		endpoints.GetCommentEndpoint,
		thread.decodeGetCommentRequest,
		thread.encodeResponse,
	)

	http.Handle("/postComment", postCommentHandler)
	http.Handle("/getComment", getCommentHandler)

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}

