// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestElasticSanName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// 23 chars
			input:    "abcdefghijklmnopqrstuvw",
			expected: true,
		},
		{
			// 24 chars
			input:    "abcdefghijklmnopqrstuvwx",
			expected: true,
		},
		{
			// 25 chars
			input:    "abcdefghijklmnopqrstuvwxy",
			expected: false,
		},
		{
			// may contain alphanumerics, dashes and underscores
			input:    "hello_world7-goodbye",
			expected: true,
		},
		{
			// must begin with an alphanumeric
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: false,
		},
		{
			// can end with an underscore
			input:    "hello_",
			expected: false,
		},
		{
			// cannot have consecutive underscore
			input:    "hello__world",
			expected: false,
		},
		{
			// cannot have consecutive dash
			input:    "hello--world",
			expected: false,
		},
		{
			// cannot have consecutive underscore or dash
			input:    "hello-_world",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// start with a number
			input:    "0abc",
			expected: true,
		},
		{
			// contain only numbers
			input:    "12345",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ElasticSanName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			if len(errors) > 0 {
				t.Logf("[DEBUG] Errors: %v", errors)
			}
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
