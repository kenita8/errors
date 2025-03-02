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
)

var (
	As = errors.As
	Is = errors.Is
)

type ConstError struct {
	error
	details []string
}

func New(text string) *ConstError {
	return &ConstError{error: errors.New(text)}
}

func (c *ConstError) Error() string {
	text := c.error.Error()
	if c.details != nil {
		text += "("
		text += strings.Join(c.details, ", ")
		text += ")"
	}
	return text
}

func (c *ConstError) Wrap(err error) *WrappedError {
	return &WrappedError{wrapper: c, wrapped: err}
}

func (c *ConstError) Details(kv ...any) *ConstError {
	c.details = []string{}
	for i := 0; i < len(kv); i += 2 {
		key := kv[i]
		var value any
		if i+1 < len(kv) {
			value = kv[i+1]
		}
		c.details = append(c.details, fmt.Sprintf("%v=%v", key, value))
	}
	return c
}

type WrappedError struct {
	wrapper error
	wrapped error
}

func (we *WrappedError) Error() string {
	if we.wrapped != nil {
		return we.wrapper.Error() + ": " + we.wrapped.Error()
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
	return we == target || errors.Is(we.wrapper, target) || errors.Is(we.wrapped, target)
}
