package repo

import (
	"fmt"
	"os"
	"testing"
)

func TestRepo(t *testing.T) {
	// Create a temporary database file for testing.
	dbPath := "test.db"
	defer os.Remove(dbPath)

	// Initialize the repository.
	r, err := New(dbPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	// Table-driven tests for Get.
	t.Run("Get", func(t *testing.T) {
		// defer r.Delete("key1")
		tests := []struct {
			key      string
			want     string
			wantErr  bool
			preamble func(*Repo) error
		}{
			{key: "key1", want: "value1", preamble: func(r *Repo) error { return r.Create("key1", "value1") }},
			{key: "key2", wantErr: true}, // Key not found
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
				if tt.preamble != nil {
					if err := tt.preamble(r); err != nil {
						t.Fatal(err)
					}
				}
				got, err := r.Get(tt.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Table-driven tests for Create.
	t.Run("Create", func(t *testing.T) {
		tests := []struct {
			key     string
			value   string
			wantErr bool
		}{
			{key: "key3", value: "value3"},
			{key: "key3", value: "value3", wantErr: true}, // Duplicate key
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
				err := r.Create(tt.key, tt.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Table-driven tests for Update.
	t.Run("Update", func(t *testing.T) {
		tests := []struct {
			key      string
			value    string
			wantErr  bool
			preamble func(*Repo) error
		}{
			{key: "key10", value: "new", preamble: func(r *Repo) error { return r.Create("key10", "old") }},
			{key: "key11", value: "value11", wantErr: true}, // Key not found
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
				if tt.preamble != nil {
					if err := tt.preamble(r); err != nil {
						t.Fatal("preamble:", err)
					}
				}
				err := r.Update(tt.key, tt.value)
				if (err != nil) != tt.wantErr {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})

	// Table-driven tests for Delete.
	t.Run("Delete", func(t *testing.T) {
		tests := []struct {
			key      string
			wantErr  bool
			preamble func(*Repo) error
		}{
			{key: "key20", preamble: func(r *Repo) error { return r.Create("key20", "value20") }},
			{key: "key21", wantErr: true}, // Deleting a non-existent key should not return an error
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
				if tt.preamble != nil {
					if err := tt.preamble(r); err != nil {
						t.Fatal(err)
					}
				}
				err := r.Delete(tt.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}
