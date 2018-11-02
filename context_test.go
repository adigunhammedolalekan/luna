package luna

import (
	"testing"
)

func TestExtractParams(t *testing.T) {

	tpl := "/hello/{name}"
	path := "/hello/luna"

	vars, err := ExtractParams(tpl, path)
	if err != nil {
		t.Error(err.Error())
	}

	if name, ok := vars["name"] . (string); !ok || name != "luna" {
		t.Error("Expected named param {name} to be luna")
	}
}

func TestMatch(t *testing.T) {

	paths := []string {"/hello/lekan/go/hammed", "/hello/12/get", "/hello/lekan/get/12"}
	regs := []string {"/hello/{name}/go/{param}", "/hello/{id}/get", "/hello/{name}/get/{id}"}

	for i, v := range paths {

		if !Match(regs[i], v) {
			t.Error("Failed to match route")
		}
	}
}
