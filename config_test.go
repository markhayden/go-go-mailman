package gogomailer

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestMain(t *testing.T) {
	assert.Equal(t, Conf["test"]["amihere"] == "yep", true)
}
