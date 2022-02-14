package main

import (
	"flag"
	"fmt"
	"github.com/TomaszDomagala/microservice-imageboard/thread"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

const (
	host     = "thread_db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "thread"
)

func main() {
	listen := flag.String("listen", ":80", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	svc := thread.NewPostgresService(fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
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

	getChildrenHandler := httptransport.NewServer(
		endpoints.GetChildrenEndpoint,
		thread.DecodeGetChildrenRequest,
		thread.EncodeResponse,
	)

	deleteThreadHandler := httptransport.NewServer(
		endpoints.DeleteThreadEndpoint,
		thread.DecodeDeleteThreadRequest,
		thread.EncodeResponse,
	)

	createThreadHandler := httptransport.NewServer(
		endpoints.CreateThreadEndpoint,
		thread.DecodeCreateThreadRequest,
		thread.EncodeResponse,
	)

	http.Handle("/postComment", postCommentHandler)
	http.Handle("/getComment", getCommentHandler)
	http.Handle("/getChildren", getChildrenHandler)
	http.Handle("/deleteThread", deleteThreadHandler)
	http.Handle("/createThread", createThreadHandler)

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
