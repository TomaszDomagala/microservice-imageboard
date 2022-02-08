package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

func main() {
	svc := NewInmemService()
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

	http.ListenAndServe(":8080", nil)
}
