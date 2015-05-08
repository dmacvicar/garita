package token

import (
	"testing"
	"reflect"
)

func TestParse(t *testing.T) {

	scope, err := ParseScope("repository:duncan/busybox:pull,push")

	if err != nil {
		t.Fail()
	}

	if (scope.Name != "duncan/busybox") {
		t.Errorf(scope.Name)
	}

	if (scope.Type != "repository") {
		t.Errorf(scope.Type)
	}

	if (scope.Namespace != "duncan") {
		t.Errorf(scope.Name)
	}

	if !reflect.DeepEqual(scope.Actions, []string{"pull", "push"}) {
		t.Errorf("%v", scope.Actions)

	}
}

func TestParseBroken(t *testing.T) {

	_, err := ParseScope("duncan/busybox:pull,push")

	if err == nil {
		t.Fail()
	}

}
