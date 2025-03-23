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
		assert.Equal(t, true, errors.Is(err, ErrNotFound))
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, false, errors.Is(err, ErrOriginal1))
		assert.Equal(t, false, errors.Is(Err1, Err2))
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
		err1 := func1()
		errA := funcA()
		assert.Equal(t, true, errors.Is(err1, Err1))
		assert.Equal(t, true, errors.Is(err1, Err2))
		assert.Equal(t, true, errors.Is(err1, Err3))
		assert.Equal(t, false, errors.Is(err1, ErrNotFound))
		assert.Equal(t, "err 1(k11=v11, k12=v12, k13=v13): err 2(k21=v21, k22=v22, k23=v23): err 3(k31=v31, k32=v32, k33=v33): error original1", err1.Error())
		assert.Equal(t, "err 1(kD1=vD1, kD2=vD2, kD3=vD3): err 2(kD1=vD1, kD2=vD2, kD3=vD3): err 3(kD1=vD1, kD2=vD2, kD3=vD3): error originalA", errA.Error())
	})
	t.Run("4: Nested Wraps", func(t *testing.T) {
		err := Err1.Wrap(Err2.Wrap(Err3))
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, true, errors.Is(err, Err2))
		assert.Equal(t, true, errors.Is(err, Err3))
		assert.Equal(t, false, errors.Is(err, ErrNotFound))
		assert.Equal(t, "err 1: err 2: err 3", err.Error())
	})

	t.Run("5: Multiple Wraps with Different Error Types", func(t *testing.T) {
		err := Err1.Wrap(ErrNotFound).Wrap(ErrOriginal1)
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, true, errors.Is(err, ErrNotFound))
		assert.Equal(t, true, errors.Is(err, ErrOriginal1))
		assert.Equal(t, false, errors.Is(err, Err2))
		assert.Equal(t, false, errors.Is(err, Err3))
		assert.Equal(t, "err 1: not found: error original1", err.Error())
	})

	t.Run("6: Wrap with No Inner Error", func(t *testing.T) {
		err := Err1.Wrap(nil)
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, false, errors.Is(err, ErrNotFound))
		assert.Equal(t, "err 1", err.Error())
	})

	t.Run("7: Wrap with Wrapped Error", func(t *testing.T) {
		innerErr := errors.New("inner wrapped error")
		err := Err1.Wrap(innerErr)
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, true, errors.Is(err, innerErr))
		assert.Equal(t, "err 1: inner wrapped error", err.Error())
	})

	t.Run("8: Wrap with Multiple Inner Errors", func(t *testing.T) {
		innerErr1 := errors.New("inner error 1")
		innerErr2 := errors.New("inner error 2")
		err := Err1.Wrap(innerErr1).Wrap(innerErr2)
		assert.Equal(t, true, errors.Is(err, Err1))
		assert.Equal(t, true, errors.Is(err, innerErr1))
		assert.Equal(t, true, errors.Is(err, innerErr2))
		assert.Equal(t, "err 1: inner error 1: inner error 2", err.Error())
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
	return ErrOriginal1
}

func funcA() error {
	err := funcB()
	if err != nil {
		return Err1.Details("kD1", "vD1", "kD2", "vD2", "kD3", "vD3").Wrap(err)
	}
	return nil
}

func funcB() error {
	err := funcC()
	if err != nil {
		return Err2.Details("kD1", "vD1", "kD2", "vD2", "kD3", "vD3").Wrap(err)
	}
	return nil
}

func funcC() error {
	err := funcD()
	if err != nil {
		return Err3.Details("kD1", "vD1", "kD2", "vD2", "kD3", "vD3").Wrap(err)
	}
	return nil
}

func funcD() error {
	return ErrOriginalA
}
