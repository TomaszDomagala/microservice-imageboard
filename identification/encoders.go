package identification

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeIdentifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request identifyRequest
	request.Ip = r.RemoteAddr
	return request, nil
}

func encodeIdentifyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

