package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestState_ReadValue(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		repoResponse string
		want         string
		wantErr      bool
	}{
		{
			name:         "ok",
			key:          "key1",
			repoResponse: "value1",
			want:         "value1",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			m := &mockRepoStorage{}
			m.On("Get", tt.key).Return(tt.repoResponse, nil)
			s := State{db: m}
			got, err := s.ReadValue(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadValue() got = %v, want %v", got, tt.want)
			}
		})
	}

}

type mockRepoStorage struct {
	mock.Mock
}

func (m *mockRepoStorage) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *mockRepoStorage) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *mockRepoStorage) Create(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *mockRepoStorage) Update(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}
