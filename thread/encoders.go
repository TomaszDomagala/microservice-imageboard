package thread

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeGetCommentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodePostCommentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request postCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
