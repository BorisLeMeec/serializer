package Serializer_test

import (
	"encoding/json"
	"testing"

	serializer "github.com/BorisLeMeec/Serializer"
	"github.com/stretchr/testify/assert"
)

type AnonymousStruct struct {
	Foo string `serializer:"public"`
}

type AnonymousString string

func TestSerialize(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		out := serializer.Serialize(nil, "")
		assert.Equal(t, nil, out)
	})

	t.Run("empty string", func(t *testing.T) {
		out := serializer.Serialize("", "")
		assert.Equal(t, "", out)
	})

	t.Run("pointer to  string", func(t *testing.T) {
		str := "foo"
		out := serializer.Serialize(&str, "")
		assert.Equal(t, &str, out)
	})

	t.Run("int", func(t *testing.T) {
		out := serializer.Serialize(0, "")
		assert.Equal(t, int64(0), out)
	})

	t.Run("bool", func(t *testing.T) {
		out := serializer.Serialize(true, "")
		assert.Equal(t, true, out)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var foo *bool
		out := serializer.Serialize(foo, "")

		var test *bool = nil
		assert.Equal(t, test, out)
	})

	t.Run("bool pointer", func(t *testing.T) {
		bar := true
		foo := &bar
		out := serializer.Serialize(foo, "")
		assert.Equal(t, &bar, out)
	})

	t.Run("struct pointer", func(t *testing.T) {
		bar := struct {
			Bar  string `serializer:"public"`
			Bar2 string
		}{Bar: "foo"}
		str, _ := json.Marshal(serializer.Serialize(&bar, "public"))
		strTest, _ := json.Marshal(&struct {
			Bar string `serializer:"public"`
		}{Bar: "foo"})
		assert.Equal(t, str, strTest)
	})

	t.Run("empty struct", func(t *testing.T) {
		out := serializer.Serialize(struct{}{}, "")
		assert.Equal(t, nil, out)
	})

	t.Run("simple struct no tags", func(t *testing.T) {
		out := serializer.Serialize(struct {
			Test string
		}{Test: "foo"}, "public")
		assert.Equal(t, nil, out)
	})

	t.Run("simple struct no serializer tags", func(t *testing.T) {
		out := serializer.Serialize(struct {
			Test string `json:"test"`
		}{Test: "foo"}, "public")
		assert.Equal(t, nil, out)
	})

	t.Run("simple struct no correct tags", func(t *testing.T) {
		out := serializer.Serialize(struct {
			Test string `serializer:"foo"`
		}{Test: "foo"}, "public")
		assert.Equal(t, nil, out)
	})

	t.Run("simple struct correct tag", func(t *testing.T) {
		out := serializer.Serialize(struct {
			Test string `serializer:"public"`
		}{Test: "foo"}, "public")
		assert.Equal(t, struct {
			Test string `serializer:"public"`
		}{Test: "foo"}, out)
	})

	t.Run("array string", func(t *testing.T) {
		out := serializer.Serialize([]string{"foo", "bar"}, "public")
		assert.Equal(t, []interface{}{"foo", "bar"}, out)
	})

	t.Run("array simple struct correct tag", func(t *testing.T) {
		out := serializer.Serialize([]struct {
			Test string `serializer:"public"`
		}{{Test: "foo"}}, "public")
		assert.Equal(t, []interface{}{struct {
			Test string `serializer:"public"`
		}{Test: "foo"}}, out)
	})

	t.Run("struct with AnonymousStruct", func(t *testing.T) {
		str, _ := json.Marshal(serializer.Serialize(struct {
			AnonymousStruct `serializer:"public"`
		}{AnonymousStruct{Foo: "foo"}}, "public"))
		strTest, _ := json.Marshal(struct {
			AnonymousStruct `serializer:"public"`
		}{AnonymousStruct{Foo: "foo"}})
		assert.Equal(t, str, strTest)
	})

	t.Run("struct with AnonymousStruct jsoned", func(t *testing.T) {
		str, _:= json.Marshal(serializer.Serialize(struct {
			AnonymousStruct `json:"foo" serializer:"public"`
		}{AnonymousStruct{Foo: "foo"}}, "public"))
		strTest, _ := json.Marshal(struct {
			AnonymousStruct AnonymousStruct `json:"foo" serializer:"public"`
		}{AnonymousStruct{
			Foo: "foo",
		}})
		assert.Equal(t, str, strTest)

	})

	t.Run("struct with AnonymousString", func(t *testing.T) {
		str, _ := json.Marshal(serializer.Serialize(struct {
			AnonymousString `json:"foo" serializer:"public"`
		}{"foo"}, "public"))
		strTest, _ := json.Marshal(struct {
			AnonymousString AnonymousString `json:"foo" serializer:"public"`
		}{"foo"})

		assert.Equal(t, str, strTest)
	})
}
