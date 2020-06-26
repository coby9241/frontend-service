package users_test

import (
	"testing"

	. "github.com/coby9241/frontend-service/internal/models/users"
	"github.com/stretchr/testify/assert"
)

func TestRoleTableName(t *testing.T) {
	role := &Role{}
	assert.Equal(t, "roles", role.TableName())
}
