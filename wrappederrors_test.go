package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ErrOriginal = errors.New("error original")
	ErrNotFound = New("not found")
	Err1        = New("err 1")
	Err2        = New("err 2")
	Err3        = New("err 3")
)

func TestWrap(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound)
		assert.Equal(t, true, errors.Is(err, ErrNotFound))
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, "err 1: not found", err.Error())
	})
	t.Run("2", func(t *testing.T) {
		err := Err2.Wrap(Err3)
		err = Err1.Wrap(err)
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, true, errors.Is(err, Err2))
		assert.Equal(t, true, errors.Is(err, Err3))
		assert.Equal(t, false, errors.Is(err, ErrNotFound))
		assert.Equal(t, "err 1: err 2: err 3", err.Error())
	})
	t.Run("3", func(t *testing.T) {
		err := func1()
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, true, errors.Is(err, Err2))
		assert.Equal(t, true, errors.Is(err, Err3))
		assert.Equal(t, false, errors.Is(err, ErrNotFound))
		assert.Equal(t, "err 1(k11=v11, k12=v12, k13=v13): err 2(k21=v21, k22=v22, k23=v23): err 3(k31=v31, k32=v32, k33=v33): error original", err.Error())
	})
}

func func1() error {
	err := func2()
	if err != nil {
		return Err1.Details("k11", "v11", "k12", "v12", "k13", "v13").Wrap(err)
	}
	return nil
}

func func2() error {
	err := func3()
	if err != nil {
		return Err2.Details("k21", "v21", "k22", "v22", "k23", "v23").Wrap(err)
	}
	return nil
}

func func3() error {
	err := func4()
	if err != nil {
		return Err3.Details("k31", "v31", "k32", "v32", "k33", "v33").Wrap(err)
	}
	return nil
}

func func4() error {
	return ErrOriginal
}
