package media

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeDeleteMediaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetMediaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

const FILE_LIMIT_MB = 2

func DecodePostMediaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	err := r.ParseMultipartForm(FILE_LIMIT_MB * 1e6)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	name := r.Form.Get("name")
	reader, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	return postMediaRequest{
		name,
		reader,
	}, nil
}

func EncodeGetResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(getMediaResponse)
	if resp.Error != nil {
		// I don't know if returning json is the best way here. Maybe plain text?
		return json.NewEncoder(w).Encode(map[string]string{"error": resp.Error.Error()})

		//_, err := w.Write([]byte(resp.Error.Error()))
		//return err
	}
	w.Header().Add("Content-Type", http.DetectContentType(resp.Data))
	_, err := w.Write(resp.Data)
	if err != nil {
		return err
	}
	return nil
}

func codeFrom(err error) int {
	switch err {
	case ErrNotExist:
		return http.StatusNotFound
	case ErrFileExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
