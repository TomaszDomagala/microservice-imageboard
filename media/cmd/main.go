package main

import (
	"flag"
	"github.com/TomaszDomagala/microservice-imageboard/media"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

func main() {
	listen := flag.String("listen", ":80", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	svc := media.NewSimpleStatelessService()
	svc = media.ServiceLoggingMiddleware(logger)(svc)

	endpoints := media.MakeServiceEndpoints(svc)

	postMediaHandler := httptransport.NewServer(endpoints.PostMediaEndpoint,
		media.DecodePostMediaRequest,
		media.EncodeResponse)

	getMediaHandler := httptransport.NewServer(endpoints.GetMediaEndpoint,
		media.DecodeGetMediaRequest,
		media.EncodeGetResponse)

	deleteMediaHandler := httptransport.NewServer(endpoints.DeleteMediaEndpoint,
		media.DecodeDeleteMediaRequest,
		media.EncodeResponse)

	http.Handle("/postMedia", postMediaHandler)
	http.Handle("/getMedia", getMediaHandler)
	http.Handle("/deleteMedia", deleteMediaHandler)
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
