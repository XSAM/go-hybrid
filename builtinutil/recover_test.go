package builtinutil

import (
	"errors"
	"io/ioutil"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	// Can take nil context
	defer Recovery(nil)
	panic("testing")
}

func TestRecoveryWithContext(t *testing.T) {
	// Can take nil context
	defer RecoveryWithContext(nil, nil)
	panic("testing")
}

func TestRecovery2(t *testing.T) {
	defer func() {
		err := recover()
		Recovery(err)
		assert.NotNil(t, err)
	}()
	panic("testing")
}

func TestStack(t *testing.T) {
	result := Stack(1)
	assert.Contains(t, string(result), "go-hybrid/builtinutil/recover_test.go")

	// Let ioutil.ReadFile() return error
	monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return nil, errors.New("mock read file failed")
	})
	defer monkey.Unpatch(ioutil.ReadFile)
	result = Stack(1)
	assert.Contains(t, string(result), "go-hybrid/builtinutil/recover_test.go")
}

func TestSource(t *testing.T) {
	bs := source(nil, 0)
	assert.Equal(t, []byte("???"), bs)

	in := [][]byte{
		[]byte("Hello world."),
		[]byte("Hello world2."),
	}
	bs = source(in, 10)
	assert.Equal(t, []byte("???"), bs)

	bs = source(in, 1)
	assert.Equal(t, []byte("Hello world."), bs)
}

func TestFunction(t *testing.T) {
	bs := function(1)

	assert.Equal(t, []byte("???"), bs)
}
