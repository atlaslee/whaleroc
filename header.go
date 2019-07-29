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
	"net"
	"unsafe"
)

const (
	SIZEOF_HEADER = SIZEOF_PROTO + SIZEOF_VERSION + 4
)

type Header struct {
	proto   [SIZEOF_PROTO]byte
	version Version
	size    uint32
}

func (this *Header) Proto() string {
	return string(this.proto[:])
}

func (this *Header) SetProto(name string) {
	copy(this.proto[:], []byte(name))
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

	copy(header.proto[:], []byte(DMT_PROTO+"/"+proto))

	if ver != nil {
		header.Version().SetBytes(ver.Bytes())
	} else {
		header.Version().SetBytes(DMT_VERSION.Bytes())
	}

	header.size = uint32(len(bytes))
	return
}

func ReadHeader(conn *net.TCPConn, proto string) (header *Header, err error) {
	h := [SIZEOF_HEADER]byte{}

	n, err := conn.Read(h[:])

	if err != nil {
		return
	}

	if n != SIZEOF_HEADER {
		return nil, ERR_BUFFER_BROKEN
	}

	header = (*Header)(unsafe.Pointer(&h))
	if header.Proto() != proto {
		return nil, ERR_UNKNOWN_PROTO
	}

	return
}

func WriteHeader(conn *net.TCPConn, proto string, ver *Version, bytes []byte) (err error) {
	header := NewHeader3(proto, ver, bytes)
	n, err := conn.Write(header.Bytes())
	if err != nil {
		return
	}

	if n != SIZEOF_HEADER {
		return ERR_BUFFER_BROKEN
	}

	return
}
