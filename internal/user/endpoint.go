package user

import (
	"context"

	"github.com/zchelalo/go_microservices_meta/meta"
	"github.com/zchelalo/go_microservices_response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetRequest struct {
		Id string
	}

	GetAllRequest struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	UpdateRequest struct {
		Id        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	DeleteRequest struct {
		Id string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(service Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(service),
		Get:    makeGetEndpoint(service),
		GetAll: makeGetAllEndpoint(service, config),
		Update: makeUpdateEndpoint(service),
		Delete: makeDeleteEndpoint(service),
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		if req.FirstName == "" {
			return nil, response.BadRequest("first name is required")
		}

		if req.LastName == "" {
			return nil, response.BadRequest("last name is required")
		}

		user, err := service.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func makeGetEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)

		user, err := service.Get(ctx, req.Id)
		if err != nil {
			return nil, response.NotFound(err.Error())
		}

		return response.OK("success", user, nil), nil
	}
}

func makeGetAllEndpoint(service Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllRequest)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := service.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := service.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", users, meta), nil
	}
}

func makeUpdateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest("first name is required")
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest("last name is required")
		}

		if err := service.Update(ctx, req.Id, req.FirstName, req.LastName, req.Email, req.Phone); err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", "user updated successfully", nil), nil
	}
}

func makeDeleteEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)

		err := service.Delete(ctx, req.Id)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", "user deleted successfully", nil), nil
	}
}
