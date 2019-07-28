/*
MIT License

Copyright (c) 2019 Atlas Lee, 4859345@qq.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package dmt

import (
	"bytes"
	"testing"
)

func TestVersionBytes(t *testing.T) {
	b := NewVersion4(1, 1, 1, 1).Bytes()
	if !bytes.Equal(b, []byte{1, 1, 1, 1}) {
		t.Fatal(b)
	}
}

func TestVersionSetBytes(t *testing.T) {
	v := NewVersion()
	v.SetBytes([]byte{1, 1, 1, 1})
	if !NewVersion4(1, 1, 1, 1).Equal(v) {
		t.Fatal(v)
	}
}

func TestVersionCmp(t *testing.T) {
	v0 := NewVersion4(1, 1, 1, 1)
	v1 := NewVersion()
	v1.SetBytes([]byte{1, 0, 1, 1})

	if v0.Cmp(v1) <= 0 {
		t.Fatal(v0, v1, v0.Cmp(v1))
	}
}

func TestVersionString(t *testing.T) {
	v := NewVersion4(1, 1, 1, 1)
	if v.String() != "v1.1.1.1" {
		t.Fatal(v)
	}
}

func TestVersionSetString(t *testing.T) {
	v := NewVersion()
	v.SetString("v1.1.1.1")
	if v.Cmp(NewVersion4(1, 1, 1, 1)) != 0 {
		t.Fatal(v)
	}
}
