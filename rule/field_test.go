package rule

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
		wantErr  bool
	}{
		{name: "eqfield success", ruleName: "eqfield", value: testStruct{Field1: "test", Field2: "test"}, wantErr: false},
		{name: "eqfield failure", ruleName: "eqfield", value: testStruct{Field1: "test", Field2: "fail"}, wantErr: true},
		{name: "nefield success", ruleName: "nefield", value: testStruct{Field1: "test", Field2: "fail"}, wantErr: false},
		{name: "nefield failure", ruleName: "nefield", value: testStruct{Field1: "test", Field2: "test"}, wantErr: true},
		{name: "gtfield success", ruleName: "gtfield", value: testStruct{Field1: 10, Field2: 5}, wantErr: false},
		{name: "gtfield failure", ruleName: "gtfield", value: testStruct{Field1: 5, Field2: 10}, wantErr: true},
		{name: "ltfield success", ruleName: "ltfield", value: testStruct{Field1: 5, Field2: 10}, wantErr: false},
		{name: "ltfield failure", ruleName: "ltfield", value: testStruct{Field1: 10, Field2: 5}, wantErr: true},
		{name: "gtefield success", ruleName: "gtefield", value: testStruct{Field1: 10, Field2: 10}, wantErr: false},
		{name: "gtefield failure", ruleName: "gtefield", value: testStruct{Field1: 5, Field2: 10}, wantErr: true},
		{name: "ltefield success", ruleName: "ltefield", value: testStruct{Field1: 10, Field2: 10}, wantErr: false},
		{name: "ltefield failure", ruleName: "ltefield", value: testStruct{Field1: 15, Field2: 10}, wantErr: true},
		{name: "fieldcontains success", ruleName: "fieldcontains", value: testStruct{Field1: "hello world", Field2: "world"}, wantErr: false},
		{name: "fieldcontains failure", ruleName: "fieldcontains", value: testStruct{Field1: "hello", Field2: "world"}, wantErr: true},
		{name: "fieldexcludes success", ruleName: "fieldexcludes", value: testStruct{Field1: "hello", Field2: "world"}, wantErr: false},
		{name: "fieldexcludes failure", ruleName: "fieldexcludes", value: testStruct{Field1: "hello world", Field2: "world"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObject().
				AddField("Field1", schema.NewField().AddValidator(r.NewValidator(tt.ruleName, "Field2"))).
				AddField("Field2", schema.NewField())
			ctx := schema.NewContext(s, data.New(tt.value))
			err := s.Validate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ctx.Errors().HasErrorCode(tt.ruleName), tt.wantErr)
		})
	}
}
