package token

import (
	"errors"
	"fmt"
	"strings"
)

type Scope struct {
	Type      string
	Name      string
	Namespace string
	Actions   []string
}

func NewScope(typ string, name string, actions []string) *Scope {

	var namespace string
	if strings.Contains(name, "/") {
		namespace = strings.Split(name, "/")[0]
	}
	return &Scope{
		Type:      typ,
		Name:      name,
		Namespace: namespace,
		Actions:   actions,
	}
}

func ParseScope(scopeString string) (*Scope, error) {

	parts := strings.Split(scopeString, ":")
	if len(parts) != 3 {
		return nil, errors.New(fmt.Sprintf("invalid scope string: '%v'", scopeString))
	}

	return NewScope(parts[0], parts[1], strings.Split(parts[2], ",")), nil
}
