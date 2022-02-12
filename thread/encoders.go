package thread

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeGetCommentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodePostCommentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request postCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
