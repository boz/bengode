package bengode

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
)

type Encoder interface {
	Encode(b io.Writer) (int, error)
}

type StringEncoder struct {
	Value string
}

type IntEncoder struct {
	Value int64
}

type UintEncoder struct {
	Value uint64
}

type ListEncoder struct {
	Value reflect.Value
}

type DictEncoder struct {
	Value reflect.Value
}

func Encode(w io.Writer, val interface{}) (int, error) {
	return EncodeValue(w, reflect.ValueOf(val))
}

func EncodeValue(w io.Writer, val reflect.Value) (int, error) {
	encoder, err := GetEncoder(val)
	if err != nil {
		return 0, err
	}
	return encoder.Encode(w)
}

func GetEncoder(val reflect.Value) (Encoder, error) {
	switch val.Kind() {
	case reflect.String:
		return &StringEncoder{Value: val.String()}, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return &IntEncoder{Value: val.Int()}, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return &UintEncoder{Value: val.Uint()}, nil
	case reflect.Map:
		if val.Type().Key().Kind() != reflect.String {
			return nil, errors.New("Map keys must be a string")
		}
		return &DictEncoder{Value: val}, nil
	case reflect.Slice, reflect.Array:
		return &ListEncoder{Value: val}, nil
	default:
		return nil, errors.New("Invalid Value Type")
	}
}

func (e *StringEncoder) Encode(w io.Writer) (int, error) {
	total := 0

	n, err := fmt.Fprintf(w, "%d:", len(e.Value))
	total += n
	if err != nil {
		return total, err
	}

	n, err = fmt.Fprintf(w, "%s", e.Value)
	total += n
	return total, err
}

func (e *IntEncoder) Encode(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "i%de", e.Value)
}

func (e *UintEncoder) Encode(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "i%de", e.Value)
}

func (e *ListEncoder) Encode(w io.Writer) (int, error) {

	total, err := fmt.Fprintf(w, "l")
	if err != nil {
		return total, err
	}

	len := e.Value.Len()

	for i := 0; i < len; i++ {
		n, err := EncodeValue(w, e.Value.Index(i))
		total += n
		if err != nil {
			return total, err
		}
	}

	n, err := fmt.Fprintf(w, "e")
	total += n
	return total, err
}

func (e *DictEncoder) Encode(w io.Writer) (int, error) {
	total, err := fmt.Fprintf(w, "d")
	if err != nil {
		return total, err
	}

	var keys stringValues = e.Value.MapKeys()
	sort.Sort(keys)

	for _, k := range keys {
		n, err := Encode(w, k.String())
		total += n
		if err != nil {
			return total, err
		}

		n, err = EncodeValue(w, e.Value.MapIndex(k))
		total += n
		if err != nil {
			return total, err
		}
	}

	n, err := fmt.Fprintf(w, "e")
	total += n
	return total, err
}

type stringValues []reflect.Value

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }
