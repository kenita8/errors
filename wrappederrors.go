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
	"fmt"
	"strings"
	"sync"
)

var (
	As = errors.As
	Is = errors.Is
)

type ConstError struct {
	err     error
	details []string
	mu      sync.Mutex
}

func New(text string) *ConstError {
	return &ConstError{err: errors.New(text)}
}

func (c *ConstError) Error() string {
	text := c.err.Error()
	if c.details != nil {
		text += "("
		text += strings.Join(c.details, ", ")
		text += ")"
	}
	return text
}

func (c *ConstError) Dup() *ConstError {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := &ConstError{
		err: c.err,
	}
	if len(c.details) <= 0 {
		return err
	}
	err.details = make([]string, len(c.details))
	copy(err.details, c.details)
	return err
}

func (c *ConstError) Wrap(err error) *WrappedError {
	newErr := c.Dup()
	return &WrappedError{wrapper: newErr, wrapped: err}
}

func (c *ConstError) Details(kv ...any) *ConstError {
	newErr := c.Dup()
	newErr.details = make([]string, 0, len(kv)/2+1)
	for i := 0; i < len(kv); i += 2 {
		key := kv[i]
		var value any
		if i+1 < len(kv) {
			value = kv[i+1]
		}
		newErr.details = append(newErr.details, fmt.Sprintf("%v=%v", key, value))
	}
	return newErr
}

func (c *ConstError) Is(target error) bool {
	if t, ok := target.(*ConstError); ok {
		return errors.Is(c.err, t.err)
	}
	return errors.Is(c.err, target)
}

func (c *ConstError) Unwrap() error {
	return c.err
}

type WrappedError struct {
	wrapper error
	wrapped error
}

func (we *WrappedError) Error() string {
	if we.wrapped != nil {
		return fmt.Sprintf("%s: %s", we.wrapper.Error(), we.wrapped.Error())
	}
	return we.wrapper.Error()
}

func (we *WrappedError) Wrap(wrappederr error) error {
	return &WrappedError{
		wrapper: we,
		wrapped: wrappederr,
	}
}

func (err *WrappedError) Unwrap() error {
	return err.wrapped
}

func (we *WrappedError) Is(target error) bool {
	return isErrorRecursive(we, target)
}

func isErrorRecursive(err error, target error) bool {
	if err == nil {
		return false
	}

	if constErr, ok := err.(*ConstError); ok {
		if constTarget, ok := target.(*ConstError); ok {
			return errors.Is(constErr.err, constTarget.err)
		}
		return errors.Is(constErr.err, target)
	}

	if wrappedErr, ok := err.(*WrappedError); ok {
		if isErrorRecursive(wrappedErr.wrapper, target) {
			return true
		}
		return isErrorRecursive(wrappedErr.wrapped, target)
	}

	return errors.Is(err, target)
}
