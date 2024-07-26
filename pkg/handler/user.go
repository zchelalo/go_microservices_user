package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/zchelalo/go_microservices_response/response"
	"github.com/zchelalo/go_microservices_user/internal/user"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {
	router := http.NewServeMux()

	opts := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("POST /users", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		opts...,
	))
	router.Handle("GET /users", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		opts...,
	))
	router.Handle("GET /users/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	))
	router.Handle("PATCH /users/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		opts...,
	))
	router.Handle("DELETE /users/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		opts...,
	))

	return router
}

func decodeCreateUser(_ context.Context, request *http.Request) (interface{}, error) {
	var req user.CreateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetUser(_ context.Context, request *http.Request) (interface{}, error) {
	id := request.PathValue("id")
	req := user.GetRequest{
		Id: id,
	}

	return req, nil
}

func decodeGetAllUser(_ context.Context, request *http.Request) (interface{}, error) {
	queries := request.URL.Query()

	limit, _ := strconv.Atoi(queries.Get("limit"))
	page, _ := strconv.Atoi(queries.Get("page"))

	req := user.GetAllRequest{
		FirstName: queries.Get("first_name"),
		LastName:  queries.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}

func decodeUpdateUser(_ context.Context, request *http.Request) (interface{}, error) {
	var req user.UpdateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	id := request.PathValue("id")
	req.Id = id

	return req, nil
}

func decodeDeleteUser(_ context.Context, request *http.Request) (interface{}, error) {
	id := request.PathValue("id")
	req := user.DeleteRequest{
		Id: id,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
