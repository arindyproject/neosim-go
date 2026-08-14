package binding

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"neosim_go/internal/shared/validator"
)

// FieldValidator adalah interface opsional untuk custom type (seperti DateOnly)
// yang butuh validasi tambahan SETELAH berhasil di-unmarshal tanpa error Go biasa.
type FieldValidator interface {
	// ValidationError mengembalikan pesan error dan true jika value tidak valid.
	ValidationError() (msg string, invalid bool)
}

// BindErrors melakukan decode JSON secara field-per-field ke struct dst,
// mengumpulkan SEMUA error tipe data yang salah (string->int, string->bool,
// format tanggal salah via DateOnly, dll) dalam format []validator.ValidationError
// yang sama dengan validator.Validate, supaya mudah digabung di handler.
//
// Return: nil jika tidak ada error sama sekali.
func BindErrors(body []byte, dst interface{}) []validator.ValidationError {
	var errs []validator.ValidationError

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return []validator.ValidationError{{
			Field:   "_body",
			Message: "format JSON tidak valid: " + err.Error(),
		}}
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return []validator.ValidationError{{
			Field:   "_internal",
			Message: "dst harus pointer ke struct",
		}}
	}
	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		name := jsonTag
		if idx := strings.Index(jsonTag, ","); idx >= 0 {
			name = jsonTag[:idx]
		}
		if name == "" {
			continue
		}

		rawVal, sent := raw[name]
		if !sent {
			continue // field tidak dikirim di JSON, skip (biar tetap optional, "required" dicek validator.Validate)
		}

		fieldVal := rv.Field(i)
		if !fieldVal.CanAddr() {
			continue
		}
		fieldPtr := fieldVal.Addr().Interface()

		// coba unmarshal per-field
		if err := json.Unmarshal(rawVal, fieldPtr); err != nil {
			errs = append(errs, validator.ValidationError{
				Field:   name,
				Message: humanizeTypeError(err, rawVal),
			})
			continue
		}

		// kalau berhasil unmarshal, cek apakah field ini punya validasi tambahan
		if fv, ok := fieldPtr.(FieldValidator); ok {
			if msg, invalid := fv.ValidationError(); invalid {
				errs = append(errs, validator.ValidationError{
					Field:   name,
					Message: msg,
				})
			}
		}
	}

	return errs
}

func humanizeTypeError(err error, rawVal json.RawMessage) string {
	if te, ok := err.(*json.UnmarshalTypeError); ok {
		expected := humanTypeName(te.Type)
		got := strings.Trim(string(rawVal), `"`)
		return fmt.Sprintf("harus berupa %s, diterima '%s'", expected, got)
	}
	return "format tidak valid: " + err.Error()
}

func humanTypeName(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return "teks"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "angka"
	case reflect.Float32, reflect.Float64:
		return "angka desimal"
	case reflect.Bool:
		return "boolean (true/false)"
	default:
		return t.String()
	}
}
