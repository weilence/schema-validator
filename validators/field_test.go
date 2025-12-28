package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func Test_compareFieldValidator(t *testing.T) {
	r := NewRegistry()
	registerField(r)

	type testStruct struct {
		Field1 any
		Field2 any
	}

	tests := []struct {
		name     string
		ruleName string
		value    testStruct
		errStr   string
	}{
		{name: "eqfield success", ruleName: "eqfield", value: testStruct{Field1: "test", Field2: "test"}, errStr: ""},
		{name: "eqfield failure", ruleName: "eqfield", value: testStruct{Field1: "test", Field2: "fail"}, errStr: "eqfield"},
		{name: "nefield success", ruleName: "nefield", value: testStruct{Field1: "test", Field2: "fail"}, errStr: ""},
		{name: "nefield failure", ruleName: "nefield", value: testStruct{Field1: "test", Field2: "test"}, errStr: "nefield"},
		{name: "gtfield success", ruleName: "gtfield", value: testStruct{Field1: 10, Field2: 5}, errStr: ""},
		{name: "gtfield failure", ruleName: "gtfield", value: testStruct{Field1: 5, Field2: 10}, errStr: "gtfield"},
		{name: "ltfield success", ruleName: "ltfield", value: testStruct{Field1: 5, Field2: 10}, errStr: ""},
		{name: "ltfield failure", ruleName: "ltfield", value: testStruct{Field1: 10, Field2: 5}, errStr: "ltfield"},
		{name: "gtefield success", ruleName: "gtefield", value: testStruct{Field1: 10, Field2: 10}, errStr: ""},
		{name: "gtefield failure", ruleName: "gtefield", value: testStruct{Field1: 5, Field2: 10}, errStr: "gtefield"},
		{name: "ltefield success", ruleName: "ltefield", value: testStruct{Field1: 10, Field2: 10}, errStr: ""},
		{name: "ltefield failure", ruleName: "ltefield", value: testStruct{Field1: 15, Field2: 10}, errStr: "ltefield"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObjectSchema().
				AddField("Field1", schema.NewFieldSchema().AddValidator(r.NewValidator(tt.ruleName, "Field2"))).
				AddField("Field2", schema.NewFieldSchema())
			ctx := schema.NewContext(s, data.New(tt.value))
			err := s.Validate(ctx)
			if tt.errStr == "" {
				assert.NoError(t, err)
				return
			}

			assert.ErrorContains(t, err, tt.errStr)
		})
	}
}
