// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"beam.apache.org/playground/backend/pkg/errors"
	"fmt"
	"strings"
	"testing"
)

func TestInvalidArgumentError(t *testing.T) {
	title := "TITLE"
	message := "MESSAGE"
	expected := fmt.Sprintf("rpc error: code = InvalidArgument desc = %s: %s", title, message)

	err := errors.InvalidArgumentError(title, message)
	res := err.Error()

	if !strings.EqualFold(res, expected) {
		t.Errorf("Unexpected message. Expected: '%s', bet received: '%s'", expected, res)
	}
}

func TestNotFoundError(t *testing.T) {
	title := "TITLE"
	message := "MESSAGE"
	expected := fmt.Sprintf("rpc error: code = NotFound desc = %s: %s", title, message)

	err := errors.NotFoundError(title, message)
	res := err.Error()

	if !strings.EqualFold(res, expected) {
		t.Errorf("Unexpected message. Expected: '%s', bet received: '%s'", expected, res)
	}
}

func TestInternalError(t *testing.T) {
	title := "TITLE"
	message := "MESSAGE"
	expected := fmt.Sprintf("rpc error: code = Internal desc = %s: %s", title, message)

	err := errors.InternalError(title, message)
	res := err.Error()

	if !strings.EqualFold(res, expected) {
		t.Errorf("Unexpected message. Expected: '%s', bet received: '%s'", expected, res)
	}
}
