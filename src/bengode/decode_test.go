package bengode_test

import (
	"bengode"
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func decode(str string) (interface{}, error) {
	return bengode.Decode(bufio.NewReader(strings.NewReader(str)))
}

func TestDecodeString(t *testing.T) {
	result, err := decode("3:foo")
	assert.Nil(t, err)
	assert.Equal(t, "foo", result.(string))
}

func TestDecodeInt(t *testing.T) {
	result, err := decode("i200e")
	assert.Nil(t, err)
	assert.Equal(t, int64(200), result.(int64))
}

func TestDecodeList(t *testing.T) {
	result, err := decode("li200ee")
	assert.Nil(t, err)
	list := result.([]interface{})
	assert.Equal(t, 1, len(list))
	assert.Equal(t, int64(200), list[0].(int64))
}

func TestDecodeDict(t *testing.T) {
	result, err := decode("d3:fooi200ee")
	assert.Nil(t, err)
	dict := result.(map[string]interface{})
	assert.Equal(t, 1, len(dict))
	assert.Equal(t, int64(200), dict["foo"].(int64))
}
