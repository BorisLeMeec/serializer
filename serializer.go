package serializer

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/fatih/structtag"
)

type Serializer struct {
}

type structFieldValue struct {
	field reflect.StructField
	value interface{}
}

func New() Serializer {
	return Serializer{}
}

func Serialize(in interface{}, tag string) interface{} {
	s := Serializer{}
	return s._serialize(in, tag)
}

func (s *Serializer) _serialize(in interface{}, tag string) interface{} {
	k := reflect.ValueOf(in).Kind()
	switch {
	case simpleType(k):
		return in
	case k == reflect.Ptr:
		return s.serializePointer(in, tag)
	case k == reflect.Struct:
		return s.serializeStruct(in, tag)
	case k == reflect.Slice:
		return s.serializeSlice(in, tag)
	default:
		fmt.Println("unknown kind: " + k.String())
	}
	return nil
}

func (s *Serializer) serializePointer(in interface{}, tag string) interface{} {
	v := reflect.ValueOf(in).Elem()
	switch v.Kind() {
	case reflect.Struct:
		s := s._serialize(v.Interface(), tag)
		return &s
	default:
		return in
	}
}

func (s *Serializer) serializeSlice(in interface{}, tag string) []interface{} {
	var out []interface{}
	v := reflect.ValueOf(in)
	for i := 0; i < v.Len(); i++ {
		if e := s._serialize(v.Index(i).Interface(), tag); e != nil {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return []interface{}{}
	}
	return out
}

func (s *Serializer) serializeStruct(st interface{}, tag string) interface{} {
	fieldsValues := s.getStructFieldsValues(st, tag)
	if len(fieldsValues) == 0 {
		return nil
	}
	var f []reflect.StructField
	for _, fv := range fieldsValues {
		f = append(f, fv.field)
	}
	out := reflect.New(reflect.StructOf(f))
	initializeStruct(out.Elem(), fieldsValues, tag)
	return out.Elem().Interface()
}

func (s *Serializer) getStructFieldsValues(st interface{}, tag string) []structFieldValue {
	var fields []structFieldValue

	values := reflect.ValueOf(st)
	types := reflect.TypeOf(st)
	for i := 0; i < types.NumField(); i++ {
		f := types.Field(i)
		if !s.keepField(f, tag) || !unicode.IsUpper(rune(f.Name[0])) {
			continue
		}
		value := s._serialize(values.Field(i).Interface(), tag)
		fields = append(fields, structFieldValue{
			field: reflect.StructField{
				Name:      f.Name,
				Type:      reflect.TypeOf(value),
				Tag:       f.Tag,
				Anonymous: f.Anonymous,
				PkgPath:   f.PkgPath,
				Index:     f.Index,
				Offset:    f.Offset,
			},
			value: value,
		})

	}

	return fields
}

func initializeStruct(s reflect.Value, source []structFieldValue, tag string) {
	model := reflect.TypeOf(s.Interface())
	for i := 0; i < s.NumField(); i++ {
		s.Field(i).Set(reflect.ValueOf(source[i].value).Convert(model.Field(i).Type))
	}
}

func (s *Serializer) keepField(f reflect.StructField, condition string) bool {
	if f.Anonymous {
		return true
	}
	tags, err := structtag.Parse(string(f.Tag))
	if err != nil {
		return false
	}
	if tag, err := tags.Get("serializer"); err == nil {
		values := append([]string{tag.Name}, tag.Options...)
		for _, val := range values {
			if val == condition {
				return true
			}
		}
		return false
	}
	return false
}

func simpleType(kind reflect.Kind) bool {
	return kind == reflect.Bool ||
		kind == reflect.String ||
		kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Float32 ||
		kind == reflect.Float64
}
