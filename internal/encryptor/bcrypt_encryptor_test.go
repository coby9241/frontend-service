package encryptor_test

import (
	"testing"

	. "github.com/coby9241/frontend-service/internal/encryptor"
	"github.com/stretchr/testify/assert"
)

func TestEncryptor(t *testing.T) {
	testPassword := "password"
	testPasswordHash := "$2a$13$RbhfZCx.tY.3stm.BDt89OibFWFPW7FL8rTK0Dw/dlwBbVBJ1s43K"

	// generate hash
	_, err := GetInstance().Digest(testPassword)
	assert.NoError(t, err)

	// compare hash
	err = GetInstance().Compare(testPasswordHash, testPassword)
	assert.NoError(t, err)
}
