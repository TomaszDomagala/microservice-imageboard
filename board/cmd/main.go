package main

import (
	"flag"
	"fmt"
	"github.com/TomaszDomagala/microservice-imageboard/board"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

const (
	host     = "board_db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "board"
)

func main() {
	listen := flag.String("listen", ":8082", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	svc := board.NewPostgresService(fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))

	svc = board.ServiceLoggingMiddleware(logger)(svc)

	endpoints := board.MakeServiceEndpoints(svc)

	getThreadsHandler := httptransport.NewServer(
		endpoints.GetThreadsEndpoint,
		board.DecodeGetThreadsRequest,
		board.EncodeResponse,
	)

	createThreadHandler := httptransport.NewServer(
		endpoints.CreateThreadEndpoint,
		board.DecodeCreateThreadRequest,
		board.EncodeResponse,
	)

	deleteThreadHandler := httptransport.NewServer(
		endpoints.DeleteThreadEndpoint,
		board.DecodeDeleteThreadRequest,
		board.EncodeResponse,
	)

	http.Handle("/getThreads", getThreadsHandler)
	http.Handle("/createThread", createThreadHandler)
	http.Handle("/deleteThread", deleteThreadHandler)

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
