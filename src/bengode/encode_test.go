package bengode_test

import (
	"bengode"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func encode(val interface{}) (string, int, error) {
	buf := new(bytes.Buffer)
	n, err := bengode.Encode(buf, val)
	return buf.String(), n, err
}

func simpleEncodeAssert(t *testing.T, val interface{}, expect string) {
	result, n, err := encode(val)
	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	assert.Equal(t, len(expect), n)
}

func TestEncodeString(t *testing.T) {
	simpleEncodeAssert(t, "foo", "3:foo")
}

func TestEncodeInt(t *testing.T) {
	simpleEncodeAssert(t, 500, "i500e")
	simpleEncodeAssert(t, uint(500), "i500e")
	simpleEncodeAssert(t, int16(500), "i500e")
	simpleEncodeAssert(t, -3, "i-3e")
}

func TestEncodeList(t *testing.T) {
	val := []int{30, 2}
	simpleEncodeAssert(t, val, "li30ei2ee")
}

func TestEncodeDict(t *testing.T) {
	val := map[string]int{"foo": 5, "bar": 10}
	simpleEncodeAssert(t, val, "d3:bari10e3:fooi5ee")
}
