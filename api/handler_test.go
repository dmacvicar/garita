package api

import (
	"testing"
	"fmt"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"strings"
)

const htpasswdPath = "../vagrant/conf/htpasswd"
const keyPath = "../vagrant/conf/ca_bundle/server.key"

type tokenResp struct {
	Token string `json:"token"`
}

func TestUnauthorized (t *testing.T) {
	assert := assert.New(t)

	handler :=  NewGaritaTokenHandler(htpasswdPath, keyPath)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("http://example.com/v2/token?account=duncan&service=registry.test.lan")
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(err)

	handler.ServeHTTP(recorder, req)

	assert.Equal(401, recorder.Code)
}

func TestTokenOutput (t *testing.T) {
	assert := assert.New(t)

	handler :=  NewGaritaTokenHandler(htpasswdPath, keyPath)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("http://example.com/v2/token?account=duncan&service=registry.test.lan")
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("duncan", "garita")
	assert.Nil(err)

	handler.ServeHTTP(recorder, req)

	assert.Equal(200, recorder.Code)

	responseJson := new(tokenResp)
	err = json.Unmarshal(recorder.Body.Bytes(), responseJson)
	assert.Nil(err)

	// JWT tokens are XXX.YYY.ZZZ
	tokenParts := strings.Split(responseJson.Token, ".")
	assert.Equal(3, len(tokenParts))
}

