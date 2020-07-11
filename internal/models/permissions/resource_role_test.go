package permissions_test

import (
	"testing"

	. "github.com/coby9241/frontend-service/internal/models/permissions"
	"github.com/stretchr/testify/assert"
)

func TestResourceRoleTableName(t *testing.T) {
	resourceRole := &ResourceRole{}
	assert.Equal(t, "resource_role", resourceRole.TableName())
}
