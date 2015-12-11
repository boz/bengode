package bengode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type Decoder interface {
	Decode(b *bufio.Reader) (interface{}, error)
}

type Encoder interface {
	Encode(b *io.Writer)
}

type StringDecoder struct{}
type IntDecoder struct{}
type DictDecoder struct{}
type ListDecoder struct{}

func (decoder *StringDecoder) Decode(b *bufio.Reader) (interface{}, error) {

	line, err := b.ReadString(':')
	if err != nil {
		return nil, err
	}

	len, err := strconv.ParseUint(line[0:len(line)-1], 10, 64)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len, len)

	_, err = b.Read(buf)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (decoder *IntDecoder) Decode(b *bufio.Reader) (interface{}, error) {
	err := consumeByte(b, 'i')
	if err != nil {
		return nil, err
	}

	line, err := b.ReadString('e')
	if err != nil {
		return nil, err
	}

	val, err := strconv.ParseInt(line[0:len(line)-1], 10, 64)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (decoder *ListDecoder) Decode(b *bufio.Reader) (interface{}, error) {
	err := consumeByte(b, 'l')
	if err != nil {
		return nil, err
	}

	vals := make([]interface{}, 0)

	for {
		chr, err := peekByte(b)
		if err != nil {
			return vals, err
		}
		if chr == 'e' {
			break
		}
		val, err := Decode(b)
		if err != nil {
			return vals, err
		}
		vals = append(vals, val)
	}
	consumeByte(b, 'e')
	return vals, nil
}

func (decoder *DictDecoder) Decode(b *bufio.Reader) (interface{}, error) {
	err := consumeByte(b, 'd')
	if err != nil {
		return nil, err
	}

	vals := make(map[string]interface{})

	for {
		chr, err := peekByte(b)
		if err != nil {
			return vals, err
		}
		if chr == 'e' {
			break
		}

		decoder, err := GetDecoder(b)
		if err != nil {
			return vals, err
		}

		var key interface{}
		switch d := decoder.(type) {
		case *StringDecoder:
			key, err = d.Decode(b)
			if err != nil {
				return vals, err
			}
		default:
			return vals, errors.New("Invalid Key")
		}

		val, err := Decode(b)
		if err != nil {
			return vals, err
		}

		vals[key.(string)] = val
	}

	consumeByte(b, 'e')
	return vals, nil
}

func Decode(b *bufio.Reader) (interface{}, error) {
	decoder, err := GetDecoder(b)
	if err != nil {
		return nil, err
	}
	return decoder.Decode(b)
}

func GetDecoder(b *bufio.Reader) (Decoder, error) {
	char, err := peekByte(b)
	if err != nil {
		return nil, err
	}
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return &StringDecoder{}, nil
	case 'i':
		return &IntDecoder{}, nil
	case 'l':
		return &ListDecoder{}, nil
	case 'd':
		return &DictDecoder{}, nil
	}
	return nil, errors.New("Invalid character")
}

func consumeByte(b *bufio.Reader, prefix byte) error {
	chr, err := b.ReadByte()
	if chr != prefix {
		return errors.New("Invalid")
	}
	if err != nil {
		return err
	}
	return nil
}

func peekByte(b *bufio.Reader) (byte, error) {
	bytes, err := b.Peek(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}
