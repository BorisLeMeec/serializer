package serializer

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/fatih/structtag"
)

type structFieldValue struct {
	field reflect.StructField
	value interface{}
}

func Serialize(in interface{}, tag string) interface{} {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Ptr:
		return serializePointer(in, tag)
	case reflect.Struct:
		s := serializeStruct(in, tag)
		return s
		return reflect.ValueOf(s).Convert(reflect.TypeOf(in))
	case reflect.Slice:
		return serializeSlice(v.Interface(), tag)
	default:
		fmt.Println("unknown kind: " + v.Kind().String())
	}
	return nil
}

func serializePointer(in interface{}, tag string) interface{} {
	v := reflect.ValueOf(in).Elem()
	switch v.Kind() {
	case reflect.Struct:
		s := Serialize(v.Interface(), tag)
		return &s
	default:
		return in
	}
}

func serializeSlice(in interface{}, tag string) []interface{} {
	var out []interface{}
	s := reflect.ValueOf(in)
	for i := 0; i < s.Len(); i++ {
		e := Serialize(s.Index(i).Interface(), tag)
		if e != nil {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return []interface{}{}
	}
	return out
}

func serializeStruct(s interface{}, tag string) interface{} {
	fieldsValues := getStructFieldsValues(s, tag)
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

func getStructFieldsValues(s interface{}, tag string) []structFieldValue {
	var fields []structFieldValue

	values := reflect.ValueOf(s)
	types := reflect.TypeOf(s)
	for i := 0; i < types.NumField(); i++ {
		f := types.Field(i)
		if !keepField(f, tag) || !unicode.IsUpper(rune(f.Name[0])) {
			continue
		}
		fields = append(fields, structFieldValue{
			field: reflect.StructField{
				Name:      f.Name,
				Type:      f.Type,
				Tag:       f.Tag,
				Anonymous: f.Anonymous,
				PkgPath:   f.PkgPath,
				Index:     f.Index,
				Offset:    f.Offset,
			},
			value: Serialize(values.Field(i).Interface(), tag),
		})

	}

	return fields
}

func initializeStruct(s reflect.Value, source []structFieldValue, tag string) {
	model := reflect.TypeOf(s.Interface())
	for i := 0; i < s.NumField(); i++ {
		fieldValue := s.Field(i)
		fieldModel := model.Field(i)
		fieldSourceValue := source[i].value
		fieldValue.Set(reflect.ValueOf(fieldSourceValue).Convert(fieldModel.Type))
	}
}

func keepField(f reflect.StructField, condition string) bool {
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
