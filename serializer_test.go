package serializer_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/BorisLeMeec/serializer"
)

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
		assert.Equal(t, 0, out)
	})

	t.Run("bool", func(t *testing.T) {
		out := serializer.Serialize(true, "")
		assert.Equal(t, true, out)
	})

	t.Run("time.Time", func(t *testing.T) {
		now := time.Now()
		out := serializer.Serialize(now, "")
		assert.Equal(t, now, out)
	})

	t.Run("simple custom type", func(t *testing.T) {
		type Foo string
		bar := Foo("foobar")
		out := serializer.Serialize(bar, "")
		assert.Equal(t, Foo("foobar"), out)
	})

	t.Run("simple custom type with own marshall", func(t *testing.T) {
		type Foo string
		bar := Foo("foobar")
		out := serializer.Serialize(bar, "")
		assert.Equal(t, Foo("foobar"), out)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var foo *bool
		out := serializer.Serialize(foo, "")

		var test *bool = nil
		assert.Equal(t, test, out)
	})

	t.Run("map string int", func(t *testing.T) {
		foo := map[string]int{"foo": 9}
		out := serializer.Serialize(foo, "")
		assert.Equal(t, foo, out)
	})

	t.Run("map string string", func(t *testing.T) {
		foo := map[string]string{"foo": "bar", "bar": "foo"}
		out := serializer.Serialize(foo, "")
		assert.Equal(t, foo, out)
	})

	type SimpleStruct struct {
		Foo string `json:"foo" serializer:"public"`
		Bar string `json:"bar"`
	}

	type SimpleStructWithMap struct {
		Foo      string            `json:"foo" serializer:"public"`
		FortyTwo map[string]string `json:"42" serializer:"public"`
		Bar      string            `json:"bar"`
	}

	t.Run("map string struct", func(t *testing.T) {
		foo := map[string]SimpleStruct{"foo": {
			Foo: "bar",
			Bar: "foo",
		}}
		out := serializer.Serialize(foo, "public")
		str1, _ := json.Marshal(out)
		str2, _ := json.Marshal(map[string]struct {
			Foo string `json:"foo" serializer:"public"`
		}{"foo": {
			Foo: "bar",
		}})
		assert.Equal(t, str2, str1)
	})

	t.Run("map string struct with map inside", func(t *testing.T) {
		foo := map[string]SimpleStructWithMap{"foo": {
			Foo:      "bar",
			FortyTwo: map[string]string{"42": "42"},
			Bar:      "foo",
		}}
		out := serializer.Serialize(foo, "public")
		str1, _ := json.Marshal(out)
		str2, _ := json.Marshal(map[string]struct {
			Foo      string            `json:"foo" serializer:"public"`
			FortyTwo map[string]string `json:"42" serializer:"public"`
		}{"foo": {
			Foo:      "bar",
			FortyTwo: map[string]string{"42": "42"},
		}})
		assert.Equal(t, str2, str1)
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

	type AnonymousStruct struct {
		Foo string `serializer:"public"`
	}

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
		str, _ := json.Marshal(serializer.Serialize(struct {
			AnonymousStruct `json:"foo" serializer:"public"`
		}{AnonymousStruct{Foo: "foo"}}, "public"))
		strTest, _ := json.Marshal(struct {
			AnonymousStruct `json:"foo" serializer:"public"`
		}{AnonymousStruct{Foo: "foo"}})
		assert.Equal(t, str, strTest)

	})

	type AnonymousString string

	t.Run("struct with AnonymousString", func(t *testing.T) {
		str, _ := json.Marshal(serializer.Serialize(struct {
			AnonymousString `json:"foo" serializer:"public"`
		}{"foo"}, "public"))
		strTest, _ := json.Marshal(struct {
			AnonymousString AnonymousString `json:"foo" serializer:"public"`
		}{"foo"})

		assert.Equal(t, str, strTest)
	})

	type DeepStruct struct {
		AnonymousStruct
		Foo string `serializer:"public"`
	}
	t.Run("struct with deep Anonymous struct", func(t *testing.T) {
		str, _ := json.Marshal(serializer.Serialize(struct {
			Test DeepStruct
			Foo  string
			Bar  string `serializer:"public"`
		}{DeepStruct{Foo: "foo"}, "bar", "foo"}, "public"))
		strTest, _ := json.Marshal(serializer.Serialize(struct {
			Test DeepStruct
			Bar  string `serializer:"public"`
		}{DeepStruct{Foo: "foo"}, "foo"}, "public"))

		assert.Equal(t, str, strTest)
	})

	type UUID struct {
		UUID string `gorm:"size:255;uniqueIndex" json:"uuid" serializer:"public"`
	}

	type BaseModel struct {
		UUID
		ID   uint `json:"id"`
		Name string
	}

	t.Run("complex struct with complex deep Anonymous struct", func(t *testing.T) {
		str, _ := json.Marshal(serializer.Serialize(struct {
			BaseModel
			Foo string
			Bar string `serializer:"public"`
		}{BaseModel{UUID{UUID: "uuid"}, 5, "test"}, "bar", "foo"}, "public"))
		strTest, _ := json.Marshal(serializer.Serialize(struct {
			BaseModel
			Foo string
			Bar string `serializer:"public"`
		}{BaseModel{UUID{UUID: "uuid"}, 5, "test"}, "bar", "foo"}, "public"))

		assert.Equal(t, str, strTest)
	})
}
