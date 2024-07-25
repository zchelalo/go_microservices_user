package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/zchelalo/go_microservices_user/internal/user"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {
	router := http.NewServeMux()

	router.Handle("POST /users", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
	))

	return router
}

func decodeCreateUser(_ context.Context, router *http.Request) (interface{}, error) {
	var req user.CreateRequest

	if err := json.NewDecoder(router.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
