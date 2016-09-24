package gommit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	appVersion = "0.0.1"

	assert.Equal(t, "v0.0.1", GetVersion(), "Must return app version")
}
