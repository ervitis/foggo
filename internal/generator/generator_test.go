package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	expectedTemplateBaseStr = `// Code generated by foggo; DO NOT EDIT.

package testdata
`
	expectedFOPMaximumStr = expectedTemplateBaseStr + `
type TestDataOption func(*TestData)

func NewTestData(options ...TestDataOption) *TestData {
	s := &TestData{}

	for _, option := range options {
		option(s)
	}

	return s
}

func WithA(A string) TestDataOption {
	return func(args *TestData) {
		args.A = A
	}
}

func WithB(B int) TestDataOption {
	return func(args *TestData) {
		args.B = B
	}
}
`
	expectedFOPMinimumStr = expectedTemplateBaseStr + `
type TestDataOption func(*TestData)

func NewTestData(options ...TestDataOption) *TestData {
	s := &TestData{}

	for _, option := range options {
		option(s)
	}

	return s
}
`
)

func TestGenerator_GenerateFOP(t *testing.T) {
	type fields struct {
		goimports bool
	}
	type args struct {
		pkgName    string
		structName string
		sts        []*StructField
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"nominal: maximum", fields{false}, args{"testdata", "TestData", []*StructField{{Name: "A", Type: "string", Ignore: false}, {Name: "B", Type: "int", Ignore: false}, {Name: "C", Type: "int", Ignore: true}}}, expectedFOPMaximumStr, assert.NoError},
		{"nominal: maximum with goimports", fields{true}, args{"testdata", "TestData", []*StructField{{Name: "A", Type: "string", Ignore: false}, {Name: "B", Type: "int", Ignore: false}, {Name: "C", Type: "int", Ignore: true}}}, expectedFOPMaximumStr, assert.NoError},
		{"nominal: minimum", fields{false}, args{"testdata", "TestData", []*StructField{}}, expectedFOPMinimumStr, assert.NoError},
		{"nominal: minimum with goimports", fields{true}, args{"testdata", "TestData", []*StructField{}}, expectedFOPMinimumStr, assert.NoError},
		{"non_nominal: have same name fields", fields{false}, args{"testdata", "TestData", []*StructField{{Name: "A"}, {Name: "a"}}}, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			g := &Generator{
				goimports: tt.fields.goimports,
			}
			got, err := g.GenerateFOP(tt.args.pkgName, tt.args.structName, tt.args.sts)
			if !tt.wantErr(t, err, fmt.Sprintf("Generator.GenerateFOP(%v)", tt.args)) {
				return
			}
			a.Equal(tt.want, got)
		})
	}
}

func TestGenerator_checkStructFieldFormat(t *testing.T) {
	type fields struct {
		goimports bool
	}
	type args struct {
		sts []*StructField
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"nominal: struct fields have a field", fields{false}, args{[]*StructField{{Name: "A", Type: "string", Ignore: false}}}, true},
		{"nominal: struct fields have some fields", fields{false}, args{[]*StructField{{Name: "A", Type: "string", Ignore: false}, {Name: "B", Type: "string", Ignore: false}, {Name: "C", Type: "string", Ignore: false}}}, true},
		{"non_nominal: struct fields have fields of same name", fields{false}, args{[]*StructField{{Name: "A", Type: "string", Ignore: false}, {Name: "a", Type: "string", Ignore: false}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				goimports: tt.fields.goimports,
			}
			if got := g.checkStructFieldFormat(tt.args.sts); got != tt.want {
				t.Errorf("checkStructFieldFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_format(t *testing.T) {
	type fields struct {
		goimports bool
	}
	type args struct {
		b *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"nominal_goimports_is_false", fields{false}, args{bytes.NewBufferString("")}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				goimports: tt.fields.goimports,
			}
			got, err := g.format(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("format() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeGenerator(t *testing.T) {
	tests := []struct {
		name string
		want *Generator
	}{
		{"nominal", &Generator{true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeGenerator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
