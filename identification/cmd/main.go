package main

import (
	"flag"
	"github.com/TomaszDomagala/microservice-imageboard/identification"
	"github.com/go-kit/kit/log"
	"os"
)

func main() {
	listen := flag.String("listen", ":8081", "HTTP listen address")
	flag.Parse()

	db, err := identification.NewDB()
	if err != nil {
		panic(err)
	}

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	service := identification.NewService(db)
	server, err := identification.NewServer(service, logger)

	if err != nil {
		panic(err)
	}

	err = server.ListenAndServe(*listen)
	_ = logger.Log("server stopped", err)
}
