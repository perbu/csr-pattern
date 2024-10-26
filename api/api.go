package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/perbu/csr-pattern/service"
	"log/slog"
	"net/http"
	"time"
)

// generate a strict Server from openapi.yaml:
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapigen.yaml openapi.yaml

type Server struct {
	logger     *slog.Logger
	httpServer *http.Server
	service    Service
}

// interface for the service used by the API Server
type Service interface {
	DeleteKeyValue(ctx context.Context, key string) error
	ReadValue(ctx context.Context, key string) (string, error)
	CreateKeyValue(ctx context.Context, key, value string) error
	UpdateValue(ctx context.Context, key, value string) error
}

// New creates a new API Server
func New(service Service, logger *slog.Logger) (*Server, error) {
	swagger, err := GetSwagger()
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil
	s := &Server{
		logger:  logger,
		service: service,
	}
	strictHandler := NewStrictHandler(s, nil)
	httpOptions := StdHTTPServerOptions{}
	handler := HandlerWithOptions(strictHandler, httpOptions)

	s.httpServer = &http.Server{
		Handler: handler,
		Addr:    ":8080",
	}
	return s, nil
}

// Run starts the API Server
func (s Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Debug("shutting down server")
		ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx2); err != nil {
			s.logger.Error("error shutting down server", "error", err)
		}
	}()
	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe: %w", err)
	}
	return nil
}

// implementation of the generated strict Server interface

func (s Server) DeleteKeyValue(ctx context.Context, request DeleteKeyValueRequestObject) (DeleteKeyValueResponseObject, error) {
	if err := s.service.DeleteKeyValue(ctx, request.Key); err != nil {
		// error translation, to avoid leaking implementation details:
		var resourceNotFoundError service.ResourceNotFoundError
		switch {
		case errors.As(err, &resourceNotFoundError):
			return DeleteKeyValue404Response{}, nil
		default:
			return nil, err
		}
	}
	return DeleteKeyValue204Response{}, nil
}

func (s Server) ReadValue(ctx context.Context, request ReadValueRequestObject) (ReadValueResponseObject, error) {
	value, err := s.service.ReadValue(ctx, request.Key)
	if err != nil {
		// error translation, to avoid leaking implementation details:
		var resourceNotFoundError service.ResourceNotFoundError
		switch {
		case errors.As(err, &resourceNotFoundError):
			return ReadValue404Response{}, nil
		default:
			return nil, err
		}
	}
	return ReadValue200JSONResponse(value), nil
}

func (s Server) CreateKeyValue(ctx context.Context, request CreateKeyValueRequestObject) (CreateKeyValueResponseObject, error) {
	if request.Body == nil || request.Body.Value == nil {
		return CreateKeyValue400Response{}, nil
	}
	if request.Key == "" {
		return CreateKeyValue400Response{}, nil
	}
	key := request.Key
	value := *request.Body.Value
	if err := s.service.CreateKeyValue(ctx, key, value); err != nil {
		// error translation, to avoid leaking implementation details:
		var resourceExistsError service.ResourceExistsError
		switch {
		case errors.As(err, &resourceExistsError):
			return CreateKeyValue409Response{}, nil
		default:
			return nil, err
		}
	}
	return CreateKeyValue201Response{}, nil
}

func (s Server) UpdateValue(ctx context.Context, request UpdateValueRequestObject) (UpdateValueResponseObject, error) {
	if request.Body == nil || request.Body.Value == nil {
		return UpdateValue400Response{}, nil
	}
	if request.Key == "" {
		return UpdateValue400Response{}, nil
	}
	key := request.Key
	value := *request.Body.Value
	if err := s.service.UpdateValue(ctx, key, value); err != nil {
		// error translation, to avoid leaking implementation details:
		var resourceNotFoundError service.ResourceNotFoundError
		switch {
		case errors.As(err, &resourceNotFoundError):
			return UpdateValue404Response{}, nil
		default:
			return nil, err
		}
	}
	return UpdateValue204Response{}, nil
}
