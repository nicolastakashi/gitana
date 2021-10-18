package gitmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInvalidBranch(t *testing.T) {
	r := &Repository{
		Url:    "https://github.com/nicolastakashi/gitana",
		Path:   "/tmp/gitana/test",
		Branch: "invalid",
	}

	_, err := r.Get(context.TODO())

	assert.NotNil(t, err, "Get invalid branch MUST return an error")
	assert.Equal(t, err.Error(), "reference not found")
}

func TestGetValidBranch(t *testing.T) {
	r := &Repository{
		Url:    "https://github.com/nicolastakashi/gitana",
		Path:   "/tmp/gitana/test",
		Branch: "main",
	}

	_, err := r.Get(context.TODO())

	assert.Nil(t, err, "Get valid branch MUST not return an error")
}

func TestGetAlreadyExistingRepo(t *testing.T) {
	r := &Repository{
		Url:    "https://github.com/nicolastakashi/gitana",
		Path:   "/tmp/gitana/test/pull",
		Branch: "main",
	}

	_, err := r.Get(context.TODO())

	assert.Nil(t, err, "Get valid branch MUST not return an error")

	_, err = r.Get(context.TODO())

	assert.Nil(t, err, "Get valid branch MUST not return an error")
}
