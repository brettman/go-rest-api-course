// +build e2e

package test


import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

const(
	BASE_URL = "http://localhost:8080"
)

func TestGetComments(t *testing.T){
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/comment")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

