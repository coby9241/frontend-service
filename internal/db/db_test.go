// +build integration

package db

import (
	"sync"
	"testing"

	"github.com/coby9241/frontend-service/internal/config"
)

func TestGetInstance(t *testing.T) {
	cases := []struct {
		name    string
		dburl   string
		wantErr bool
	}{
		{
			name:  "test get db success",
			dburl: config.GetInstance().DatabaseURL,
		},
		{
			name:    "test get db failure",
			dburl:   "postgres://postgres@definitelydontexist:5432/admin_dev?sslmode=disable",
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tc.wantErr {
					t.Errorf("failed test to load db. wantErr: %v, err: %v", tc.wantErr, err)
				}
			}()

			once = *new(sync.Once)
			databaseURL = tc.dburl
			_ = GetInstance()
		})
	}
}
