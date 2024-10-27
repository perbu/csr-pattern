package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/perbu/csr-pattern/repo"
	"log/slog"
)

type ResourceNotFoundError struct {
	resource string
}

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("resource not found: '%s'", e.resource)
}

func NewResourceNotFoundError(resource string) ResourceNotFoundError {
	return ResourceNotFoundError{resource: resource}
}

type ResourceExistsError struct {
	resource string
}

func (e ResourceExistsError) Error() string {
	return fmt.Sprintf("resource already exists: '%s'", e.resource)
}

func NewResourceExistsError(resource string) ResourceExistsError {
	return ResourceExistsError{resource: resource}
}

type RepoError struct {
	err error
}

func (e RepoError) Error() string {
	return fmt.Sprintf("repo error: %s", e.err)
}

func NewRepoError(err error) RepoError {
	return RepoError{err: err}
}

// RepoStorage is the interface that wraps the basic operations on the repository.
type RepoStorage interface {
	Delete(key string) error
	Get(key string) (string, error)
	Create(key, value string) error
	Update(key, value string) error
}

type State struct {
	// State is the state of the service.
	db     RepoStorage
	logger *slog.Logger
}

func New(r *repo.Repo, logger *slog.Logger) (State, error) {
	return State{db: r, logger: logger}, nil
}

func (s State) CreateKeyValue(ctx context.Context, key, value string) error {
	err := s.db.Create(key, value)
	if err == nil {
		return nil
	}
	// error translation, to avoid leaking implementation details:
	var keyExistsError repo.KeyExistsError
	switch {
	case errors.As(err, &keyExistsError):
		return NewResourceExistsError(key)
	default:
		return NewRepoError(err)
	}
}

func (s State) ReadValue(ctx context.Context, key string) (string, error) {
	value, err := s.db.Get(key)
	if err == nil {
		return value, nil
	}
	// error translation, to avoid leaking implementation details:
	var keyNotFoundError repo.KeyNotFoundError
	switch {
	case errors.As(err, &keyNotFoundError):
		return "", NewResourceNotFoundError(key)
	default:
		return "", NewRepoError(err)
	}
}

func (s State) UpdateValue(ctx context.Context, key, value string) error {
	return s.db.Update(key, value)
}

func (s State) DeleteKeyValue(ctx context.Context, key string) error {
	err := s.db.Delete(key)
	if err == nil {
		return nil
	}
	// error translation, to avoid leaking implementation details:
	var keyNotFoundError repo.KeyNotFoundError
	switch {
	case errors.As(err, &keyNotFoundError):
		return NewResourceNotFoundError(key)
	default:
		return NewRepoError(err)
	}
}
