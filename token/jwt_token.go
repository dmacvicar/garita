package token

import (
	"crypto/rsa"
	"io/ioutil"
	"time"
	"os"
	libjwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/nu7hatch/gouuid"
	utils "github.com/dmacvicar/garita/utils"
)

type JwtToken struct {
	Account string
	Service string
	Scope *Scope
	privKey *rsa.PrivateKey
}

type accessItem struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Actions []string `json:"actions"`
}

func NewJwtToken(account string, service string, scope *Scope, keyPath string) (*JwtToken, error) {
	token := new(JwtToken)
	token.Account = account
	token.Service = service
	token.Scope = scope

	if err := token.parsePrivateKey(keyPath); err != nil {
		return nil, err
	}
	return token, nil
}

func (t *JwtToken) jwtKid() (string, error) {
	if kid, err := utils.KeyIDFromCryptoKey(t.privKey.Public()); err != nil {
		return "", err
	} else {
		return kid, nil
	}
}

func (t *JwtToken) parsePrivateKey(keyPath string) error {
	if pem, err := ioutil.ReadFile(keyPath); err != nil {
		return err
	} else {
		if privKey, err := libjwt.ParseRSAPrivateKeyFromPEM(pem); err != nil {
			return err
		} else {
			t.privKey = privKey
			return nil
		}

	}
}

func (t *JwtToken) notBefore() time.Time {
	return time.Now().Add(time.Second * -5)
}

func (t *JwtToken) issuedAt() time.Time {
	return t.notBefore()
}

func (t *JwtToken) expires() time.Time {
	return time.Now().Add(time.Minute * 5)
}

func (t *JwtToken) issuer() time.Time {
	return time.Now().Add(time.Minute * 5)
}

func (t *JwtToken) jwtId() (string, error) {
	if jti, err := uuid.NewV4(); err != nil {
		return "", err
	} else {
		return jti.String(), nil
	}
}

func (t *JwtToken) singleAction() accessItem {
	action := accessItem{}
	action.Type = t.Scope.Type
	action.Name = t.Scope.Name

	// only allow push pull if scope namespace
	// is the same as the authenticated account
	if (t.Account == t.Scope.Namespace) {
		action.Actions = t.Scope.Actions
	} else {
		action.Actions = []string{}
	}
	return action
}

func (t *JwtToken) authorizedAccess() []accessItem {
	return []accessItem{ t.singleAction() }
}

func (t* JwtToken) Claim() map[string]interface{} {

	claims := make(map[string]interface{})

	fqdn, err := os.Hostname()
	if err != nil {
		claims["iss"] = "garita"
	} else {
		claims["iss"] = fqdn
	}

	claims["sub"] = t.Account
	claims["aud"] = t.Service

	claims["exp"] = t.expires().Unix()
	claims["nbf"] = t.notBefore().Unix()
	claims["iat"] = t.issuedAt().Unix()

	if id, err := t.jwtId(); err != nil {
		claims["jti"] = id
	}

	if (t.Scope != nil) {
		claims["access"] = t.authorizedAccess()
	}

	return claims
}

func (t* JwtToken) SignedString() (string, error) {

	// now create the token
	jwt := libjwt.New(libjwt.SigningMethodRS256)
	jwt.Claims = t.Claim()
	if kid, err := t.jwtKid(); err != nil {
		return "", err
	} else {
		jwt.Header["kid"] = kid
	}

	signed, err := jwt.SignedString(t.privKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

