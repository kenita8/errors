// Copyright 2024 kenita8
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ErrOriginal1 = errors.New("error original1")
	ErrOriginalA = errors.New("error originalA")
	ErrNotFound  = New("not found")
	Err1         = New("err 1")
	Err2         = New("err 2")
	Err3         = New("err 3")
)

func TestWrap(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound)
		assert.Equal(t, true, Is(err, ErrNotFound))
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, false, Is(err, ErrOriginal1))
		assert.Equal(t, false, Is(Err1, Err2))
		assert.Equal(t, "err 1: not found", err.Error())
	})
	t.Run("2", func(t *testing.T) {
		err := Err2.Wrap(Err3)
		err = Err1.Wrap(err)
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, true, Is(err, Err2))
		assert.Equal(t, true, Is(err, Err3))
		assert.Equal(t, false, Is(err, ErrNotFound))
		assert.Equal(t, "err 1: err 2: err 3", err.Error())
	})
	t.Run("3", func(t *testing.T) {
		err1 := func1()
		errA := funcA()
		assert.Equal(t, true, Is(err1, Err1))
		assert.Equal(t, true, Is(err1, Err2))
		assert.Equal(t, true, Is(err1, Err3))
		assert.Equal(t, false, Is(err1, ErrNotFound))
		assert.Equal(t, "err 1(k11=v11, k12=v12, k13=v13): err 2(k21=v21, k22=v22, k23=v23): err 3(k31=v31, k32=v32, k33=v33): error original1", err1.Error())
		assert.Equal(t, "err 1(kD1=vD1, kD2=vD2, kD3=vD3): err 2(kD1=vD1, kD2=vD2, kD3=vD3): err 3(kD1=vD1, kD2=vD2, kD3=vD3): error originalA", errA.Error())
	})
	t.Run("4: Nested Wraps", func(t *testing.T) {
		err := Err1.Wrap(Err2.Wrap(Err3))
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, true, Is(err, Err2))
		assert.Equal(t, true, Is(err, Err3))
		assert.Equal(t, false, Is(err, ErrNotFound))
		assert.Equal(t, "err 1: err 2: err 3", err.Error())
	})

	t.Run("5: Multiple Wraps with Different Error Types", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound).Wrap(ErrOriginal1)
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, true, Is(err, ErrNotFound))
		assert.Equal(t, true, Is(err, ErrOriginal1))
		assert.Equal(t, false, Is(err, Err2))
		assert.Equal(t, false, Is(err, Err3))
		assert.Equal(t, "err 1: not found: error original1", err.Error())
	})

	t.Run("6: Wrap with No Inner Error", func(t *testing.T) {
		err := Err1.Wrap(nil)
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, false, Is(err, ErrNotFound))
		assert.Equal(t, "err 1", err.Error())
	})

	t.Run("7: Wrap with Wrapped Error", func(t *testing.T) {
		innerErr := errors.New("inner wrapped error")
		err := Err1.Wrap(innerErr)
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, true, Is(err, innerErr))
		assert.Equal(t, "err 1: inner wrapped error", err.Error())
	})

	t.Run("8: Wrap with Multiple Inner Errors", func(t *testing.T) {
		innerErr1 := errors.New("inner error 1")
		innerErr2 := errors.New("inner error 2")
		err := Err1.Wrap(innerErr1).Wrap(innerErr2)
		assert.Equal(t, true, Is(err, Err1))
		assert.Equal(t, true, Is(err, innerErr1))
		assert.Equal(t, true, Is(err, innerErr2))
		assert.Equal(t, "err 1: inner error 1: inner error 2", err.Error())
	})
}

func func1() error {
	err := func2()
	if err != nil {
		return Err1.WithDetails("k11", "v11", "k12", "v12", "k13", "v13").Wrap(err)
	}
	return nil
}

func func2() error {
	err := func3()
	if err != nil {
		return Err2.WithDetails("k21", "v21", "k22", "v22", "k23", "v23").Wrap(err)
	}
	return nil
}

func func3() error {
	err := func4()
	if err != nil {
		return Err3.WithDetails("k31", "v31", "k32", "v32", "k33", "v33").Wrap(err)
	}
	return nil
}

func func4() error {
	return ErrOriginal1
}

func funcA() error {
	err := funcB()
	if err != nil {
		return Err1.WithDetails("kD1", "vD1", "kD2", "vD2", "kD3", "vD3").Wrap(err)
	}
	return nil
}

func funcB() error {
	err := funcC()
	if err != nil {
		return Err2.WithDetails("kD1", "vD1", "kD2", "vD2", "kD3", "vD3").Wrap(err)
	}
	return nil
}

func funcC() error {
	err := funcD()
	if err != nil {
		return Err3.WithDetails("kD1", "vD1", "kD2", "vD2", "kD3", "vD3").Wrap(err)
	}
	return nil
}

func funcD() error {
	return ErrOriginalA
}

func TestWithDetails(t *testing.T) {
	t.Run("No details", func(t *testing.T) {
		err := Err1.WithDetails()
		assert.Equal(t, "err 1", err.Error())
	})

	t.Run("Odd number of details", func(t *testing.T) {
		err := Err1.WithDetails("key1")
		assert.Equal(t, "err 1(key1=<nil>)", err.Error())
	})

	t.Run("Different types of details", func(t *testing.T) {
		err := Err1.WithDetails("count", 123, "enabled", true)
		assert.Equal(t, "err 1(count=123, enabled=true)", err.Error())
	})

	t.Run("WithDetails on SentinelError", func(t *testing.T) {
		err := New("sentinel").WithDetails("info", "data")
		assert.Equal(t, "sentinel(info=data)", err.Error())
	})
}

func TestAs(t *testing.T) {
	t.Run("As SentinelError", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound)
		var target *SentinelError
		assert.True(t, As(err, &target))
		assert.Equal(t, Err1, target)
	})

	t.Run("As standard error", func(t *testing.T) {
		origErr := &OrgError{text: "original1"}
		err := Err1.Wrap(origErr)
		var target *OrgError
		assert.True(t, As(err, &target))
		assert.NotNil(t, target)
		assert.Equal(t, origErr, target)
	})

	t.Run("As cause error", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound)
		var target *SentinelError
		assert.True(t, As(err, &target)) // Err1 が target に入る
		target = nil
		assert.True(t, As(errors.Unwrap(err), &target)) // ErrNotFound は SentinelError
		assert.Equal(t, ErrNotFound, target)
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("Single Wrap", func(t *testing.T) {
		wrapped := Err1.Wrap(ErrNotFound)
		unwrapped := errors.Unwrap(wrapped)
		assert.Equal(t, ErrNotFound, unwrapped)
	})

	t.Run("Multiple Wraps", func(t *testing.T) {
		wrapped := Err1.Wrap(Err2.Wrap(Err3))
		unwrapped1 := errors.Unwrap(wrapped)
		if e, ok := unwrapped1.(*WrappedError); ok {
			assert.True(t, e.Is(Err2))
		} else {
			assert.Fail(t, "fail")
		}
		unwrapped2 := errors.Unwrap(unwrapped1)
		assert.True(t, errors.Is(Err3, unwrapped2))
		unwrapped3 := errors.Unwrap(unwrapped2)
		assert.Nil(t, unwrapped3)
	})

	t.Run("No cause", func(t *testing.T) {
		err := New("no cause error")
		unwrapped := errors.Unwrap(err)
		assert.Nil(t, unwrapped)
	})
}

func TestErrorMessageFormat(t *testing.T) {
	t.Run("Details with special characters", func(t *testing.T) {
		err := Err1.WithDetails("key with space", "value with comma,", "key=equal", "value")
		assert.Contains(t, err.Error(), "key with space=value with comma,")
		assert.Contains(t, err.Error(), "key=equal=value")
	})

	t.Run("Wrap error with empty message", func(t *testing.T) {
		emptyErr := errors.New("")
		err := Err1.Wrap(emptyErr)
		assert.Equal(t, "err 1: ", err.Error()) // 末尾のコロンは実装による
	})
}

func TestIs(t *testing.T) {
	t.Run("Is with nil", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound)
		assert.False(t, Is(err, nil))
	})

	t.Run("Is with same instance", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound)
		assert.True(t, Is(err, err))
	})
}

type OrgError struct {
	text string
}

func (e *OrgError) Error() string {
	return e.text
}
