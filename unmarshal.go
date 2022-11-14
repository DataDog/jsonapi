package jsonapi

import (
	"encoding/json"
	"reflect"
)

// Unmarshaler is configured internally via UnmarshalOption's passed to Unmarshal.
// It's used to configure the Unmarshaling by decoding optional fields like Meta.
type Unmarshaler struct {
	unmarshalMeta bool
	meta          any
}

// UnmarshalOption allows for configuration of Unmarshaling.
type UnmarshalOption func(m *Unmarshaler)

// UnmarshalMeta decodes Document.Meta into the given interface when unmarshaling.
func UnmarshalMeta(meta any) UnmarshalOption {
	return func(m *Unmarshaler) {
		m.unmarshalMeta = true
		m.meta = meta
	}
}

// Unmarshal parses the json:api encoded data and stores the result in the value pointed to by v.
// If v is nil or not a pointer, Unmarshal returns an error.
func Unmarshal(data []byte, v any, opts ...UnmarshalOption) (err error) {
	defer func() {
		// because we make use of reflect we must recover any panics
		if rvr := recover(); rvr != nil {
			err = recoverError(rvr)
			return
		}
	}()

	m := new(Unmarshaler)
	for _, opt := range opts {
		opt(m)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		err = &TypeError{Actual: rv.Kind().String(), Expected: []string{"non-nil pointer"}}
		return
	}

	var d document
	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	err = d.unmarshal(v, m)

	return
}

func (d *document) unmarshal(v any, m *Unmarshaler) (err error) {
	// this means we couldn't decode anything (e.g. {}, [], ...)
	if len(d.DataMany) == 0 && d.DataOne == nil {
		err = &RequestBodyError{t: v}
		return
	}

	// verify full-linkage in-case this is a compound document
	if err = d.verifyFullLinkage(); err != nil {
		return
	}

	if d.hasMany {
		err = unmarshalResourceObjects(d.DataMany, v)
		if err != nil {
			return
		}
	} else {
		err = d.DataOne.unmarshal(v)
		if err != nil {
			return
		}
	}

	err = d.unmarshalOptionalFields(m)

	return

}

func (d *document) unmarshalOptionalFields(m *Unmarshaler) error {
	if m.unmarshalMeta {
		b, err := json.Marshal(d.Meta)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, m.meta); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalResourceObjects(ros []*resourceObject, v any) error {
	outType := derefType(reflect.TypeOf(v))
	outValue := derefValue(reflect.ValueOf(v))

	// first, it must be a struct since we'll be parsing the jsonapi struct tags
	if outType.Kind() != reflect.Slice {
		return &TypeError{Actual: outType.String(), Expected: []string{"slice"}}
	}

	for _, ro := range ros {
		// unmarshal the resource object into an empty value of the slices element type
		outElem := reflect.New(derefType(outType.Elem())).Interface()
		if err := ro.unmarshal(outElem); err != nil {
			return err
		}

		// reflect.New creates a pointer, so if our slices underlying type
		// is not a pointer we must dereference the value before appending it
		outElemValue := reflect.ValueOf(outElem)
		if outType.Elem().Kind() != reflect.Pointer {
			outElemValue = derefValue(outElemValue)
		}

		// append the unmarshaled resource object to the result slice
		outValue = reflect.Append(outValue, outElemValue)
	}

	// set the value of the passed in object to our result
	reflect.ValueOf(v).Elem().Set(outValue)

	return nil
}

func (ro *resourceObject) unmarshal(v any) error {
	// first, it must be a struct since we'll be parsing the jsonapi struct tags
	vt := reflect.TypeOf(v)
	if derefType(vt).Kind() != reflect.Struct {
		return &TypeError{Actual: vt.String(), Expected: []string{"struct"}}
	}

	if err := ro.unmarshalPrimary(v); err != nil {
		return err
	}

	if err := ro.unmarshalAttributes(v); err != nil {
		return err
	}

	if err := ro.unmarshalMeta(v); err != nil {
		return err
	}

	return nil
}

func (ro *resourceObject) unmarshalPrimary(v any) error {
	setPrimary := false
	rv := derefValue(reflect.ValueOf(v))
	rt := reflect.TypeOf(rv.Interface())

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		tag, err := parseJSONAPITag(ft)
		if err != nil {
			return err
		}
		if tag == nil || tag.directive != primary {
			continue
		}
		if setPrimary {
			return ErrUnmarshalDuplicatePrimaryField
		}

		if ro.Type != tag.resourceType {
			return &TypeError{Actual: ro.Type, Expected: []string{tag.resourceType}}
		}

		// to unmarshal the id we follow these rules
		//     1. Use UnmarshalIdentifier if it is implemented
		//     2. Use the value directly if it is a string
		//     3. Fail

		if vu, ok := v.(UnmarshalIdentifier); ok {
			if err := vu.UnmarshalID(ro.ID); err != nil {
				return err
			}
			setPrimary = true
			continue
		}

		if fv.Kind() == reflect.String {
			fv.SetString(ro.ID)
			setPrimary = true
			continue
		}

		return ErrUnmarshalInvalidPrimaryField
	}

	return nil
}

func (ro *resourceObject) unmarshalAttributes(v any) error {
	b, err := json.Marshal(ro.Attributes)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

func (ro *resourceObject) unmarshalMeta(v any) error {
	if ro.Meta == nil {
		return nil
	}

	rv := derefValue(reflect.ValueOf(v))
	rt := reflect.TypeOf(rv.Interface())
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)

		tag, err := parseJSONAPITag(rt.Field(i))
		if err != nil {
			return err
		}
		if tag == nil || tag.directive != meta {
			continue
		}

		b, err := json.Marshal(ro.Meta)
		if err != nil {
			return err
		}

		v := reflect.New(derefType(f.Type())).Interface()
		err = json.Unmarshal(b, v)
		if err != nil {
			return err
		}

		// reflect.New creates a pointer so if the fields type is
		// not a pointer, dereference the unmarshaled value
		vv := reflect.ValueOf(v)
		if f.Kind() != reflect.Pointer {
			vv = derefValue(vv)
		}

		f.Set(vv)
	}

	return nil
}
