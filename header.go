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
	"unsafe"
)

const (
	SIZEOF_PROTO  = 8
	SIZEOF_HEADER = SIZEOF_PROTO + SIZEOF_VERSION + 4
	MAINPROTO     = "WRCP" // WhaleRoc Chain Protocol
)

var (
	PROTOVER = NewVersion4(0, 1, 0, 0)
)

type Header struct {
	name    [SIZEOF_PROTO]byte
	version Version
	size    uint32
}

func (this *Header) Name() string {
	return string(this.name[:])
}

func (this *Header) SetName(name string) {
	copy(this.name[:], []byte(name))
}

func (this *Header) Version() *Version {
	return &this.version
}

func (this *Header) SetVersion(v *Version) {
	this.Version().SetBytes(v.Bytes())
}

func (this *Header) Size() int {
	return int(this.size)
}

func (this *Header) SetSize(size int) {
	this.size = uint32(size)
}

func (this *Header) Bytes() []byte {
	return (*[SIZEOF_HEADER]byte)(unsafe.Pointer(this))[:]
}

func (this *Header) SetBytes(bytes []byte) *Header {
	copy(this.Bytes(), bytes)
	return this
}

func NewHeader() *Header {
	return &Header{}
}

func NewHeader1(bytes []byte) *Header {
	return NewHeader3("", nil, bytes)
}

func NewHeader2(proto string, bytes []byte) *Header {
	return NewHeader3(proto, nil, bytes)
}

func NewHeader3(proto string, ver *Version, bytes []byte) (header *Header) {
	header = NewHeader()

	copy(header.name[:], []byte(MAINPROTO+"/"+proto))

	if ver != nil {
		header.Version().SetBytes(ver.Bytes())
	} else {
		header.Version().SetBytes(PROTOVER.Bytes())
	}

	header.size = uint32(len(bytes))
	return
}
