package token

import (
	"log"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	utils "github.com/dmacvicar/garita/utils"
)

func TestJwtTokenProperties(t *testing.T) {
	assert := assert.New(t)

	const keyPath = "../vagrant/conf/ca_bundle/server.key"

	scope := NewScope("repository", "duncan", []string{"push", "pull"})

	token, err :=NewJwtToken("duncan", "registry.test.lan", scope, keyPath)
	assert.Nil(err)

	kid, err := token.jwtKid()
	assert.Nil(err)

	assert.Equal("PTWT:FNJE:7TW7:ULI7:DZQA:JJJI:RDJQ:2M76:HD6G:ZRSC:VPIF:O5BU", kid)

	log.Printf(utils.PrettyPrint(token.Claim()))

	assert.Equal("registry.test.lan", token.Claim()["aud"])

	signed, err := token.SignedString()
	assert.Nil(err)

	tokenParts := strings.Split(signed, ".")
	assert.Equal(3, len(tokenParts))
}

