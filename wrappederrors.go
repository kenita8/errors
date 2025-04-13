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

type SentinelError struct {
	error
}

func New(msg string) *SentinelError {
	return &SentinelError{
		error: errors.New(msg),
	}
}

func (s *SentinelError) Error() string {
	return s.error.Error()
}

func (s *SentinelError) Wrap(err error) *WrappedError {
	return &WrappedError{
		err:     s,
		details: nil,
		cause:   err,
	}
}

func (s *SentinelError) WithDetails(kv ...any) *WrappedError {
	return &WrappedError{
		err:     s,
		details: kv,
		cause:   nil,
	}
}

func (s *SentinelError) Is(target error) bool {
	return errors.Is(s.error, target)
}

type WrappedError struct {
	err     error
	details []any
	cause   error
}

func (s *WrappedError) Error() string {
	var sb strings.Builder
	sb.WriteString(s.err.Error())
	if len(s.details) > 0 {
		keyValues := make([]string, 0, len(s.details)/2+1)
		for i := 0; i < len(s.details); i += 2 {
			key := s.details[i]
			var value any
			if i+1 < len(s.details) {
				value = s.details[i+1]
			}
			keyValues = append(keyValues, fmt.Sprintf("%v=%v", key, value))
		}
		sb.WriteString("(")
		sb.WriteString(strings.Join(keyValues, ", "))
		sb.WriteString(")")
	}
	if s.cause != nil {
		sb.WriteString(": ")
		sb.WriteString(s.cause.Error())
	}
	return sb.String()
}

func (s *WrappedError) Wrap(err error) *WrappedError {
	return &WrappedError{
		err:     s,
		details: nil,
		cause:   err,
	}
}

func (s *WrappedError) WithDetails(kv ...any) *WrappedError {
	newError := &WrappedError{
		err:     s.err,
		details: make([]any, len(kv)),
		cause:   s.cause,
	}
	copy(newError.details, kv)
	return newError
}

func (e *WrappedError) Is(target error) bool {
	return errors.Is(e.err, target) || errors.Is(e.cause, target)
}

func (e *WrappedError) As(target any) bool {
	return errors.As(e.err, target) || errors.As(e.cause, target)
}

func (e *WrappedError) Unwrap() error {
	return e.cause
}
