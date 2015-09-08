package token

import (
	utils "github.com/dmacvicar/garita/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

func TestJwtTokenProperties(t *testing.T) {
	assert := assert.New(t)

	const keyPath = "../vagrant/conf/ca_bundle/server.key"

	scope := NewScope("repository", "duncan", []string{"push", "pull"})

	token, err := NewJwtToken("duncan", "registry.test.lan", scope, keyPath)
	assert.Nil(err)

	kid, err := token.jwtKid()
	assert.Nil(err)

	assert.Equal("NSN7:VDFR:FTW6:WBBB:7WQK:ABNJ:7CI5:M6YU:7FSD:QS45:A2BR:PAMO", kid)

	log.Printf(utils.PrettyPrint(token.Claim()))

	assert.Equal("registry.test.lan", token.Claim()["aud"])

	signed, err := token.SignedString()
	assert.Nil(err)

	tokenParts := strings.Split(signed, ".")
	assert.Equal(3, len(tokenParts))
}
