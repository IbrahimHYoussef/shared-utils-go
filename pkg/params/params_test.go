package params_test

import (
	"net/http/httptest"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/params"
)

type TestQueryParams struct {
	Int               int     `json:"int"`
	String            string  `json:"string"`
	Float64           float64 `json:"float64"`
	Float32           float32 `json:"float32"`
	Bool              bool    `json:"bool"`
	IntNoneParsed     int
	StringNoneParsed  string
	Float64NoneParsed float64
	Float32NoneParsed float32
}

func TestBindQuery(t *testing.T) {
	r := httptest.NewRequest("GET", "/p?int=1&string=name&float64=3.5&float32=2.5&bool=false", nil)

	query, err := params.BindQuery[TestQueryParams](r)
	if err != nil {
		t.Errorf("failed To Parse with error %s", err)
		return
	}
	if query == nil {
		t.Errorf("Error query is nil and the error is nil")
		return
	}

	wantStruct := TestQueryParams{
		Int:     1,
		String:  "name",
		Float64: 3.5,
		Float32: 2.5,
		Bool:    false,
	}

	if *query != wantStruct {
		t.Error("The Test Param and the wanted Param Do not match")
		return
	}
}
