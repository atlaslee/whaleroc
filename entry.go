/*
The MIT License (MIT)

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
	PROTO_ENTRY        = "NTR"
	SIZEOF_ENTRY_PROTO = len(PROTO_BLOCK)
	SIZEOF_ENTRY       = SIZEOF_ENTRY_PROTO + SIZEOF_VERSION + SIZEOF_HASH + SIZEOF_ADDRESS + 8 + SIZEOF_HASH + 4 + 1 + SIZEOF_SIGNATURE
)

// 条目接口
// 该接口不得修改，以便向前兼容
//
// 用于保存传播内容
type EntryI interface {
	Proto() string // 用于反序列化校验
	SetProto(string)
	Version() *Version // 用于指定处理的版本
	Id() [SIZEOF_HASH]byte
	From() [SIZEOF_ADDRESS]byte
	Timestamp() uint64
	Class() [SIZEOF_HASH]byte // 数据的类型，决定处理接口
	Size() uint32
	EncryptType() byte
	Signature() []byte
	Data() []byte
	Bytes() []byte
	SetBytes([]byte) (EntryI, error)
}

// 条目
// 0_1为版本号。发布过的结构不得改变,只能发新结构
// 发布过的属性不得修改，不得调整已有顺序，只能在signature后新增
type Entry_0_1 struct {
	proto       [SIZEOF_BLOCK_PROTO]byte
	version     Version
	id          [SIZEOF_HASH]byte
	from        [SIZEOF_ADDRESS]byte
	timestamp   uint64
	class       [SIZEOF_HASH]byte
	size        uint32
	encryptType byte
	signature   [SIZEOF_SIGNATURE]byte
	data        []byte
}

func (this *Entry_0_1) Proto() string {
	return string(this.proto[:])
}

func (this *Entry_0_1) SetProto(proto string) {
	copy(this.proto[:], []byte(proto))
}

func (this *Entry_0_1) Version() *Version {
	return &this.version
}

func (this *Entry_0_1) Id() [SIZEOF_HASH]byte {
	return this.id
}

func (this *Entry_0_1) From() [SIZEOF_ADDRESS]byte {
	return this.from
}

func (this *Entry_0_1) Timestamp() uint64 {
	return this.timestamp
}

func (this *Entry_0_1) Class() [SIZEOF_HASH]byte {
	return this.class
}

func (this *Entry_0_1) Size() uint32 {
	return this.size
}

func (this *Entry_0_1) EncryptType() byte {
	return this.encryptType
}

func (this *Entry_0_1) Signature() []byte {
	return this.signature[:]
}

func (this *Entry_0_1) Data() []byte {
	return this.data
}

func (this *Entry_0_1) Bytes() (b []byte) {
	b = (*(*[SIZEOF_ENTRY]byte)(unsafe.Pointer(this)))[:]
	if this.size == 0 {
		return
	}

	b = append(b, this.data...)
	return
}

func (this *Entry_0_1) SetBytes(b []byte) (entry EntryI, err error) {
	copy((*(*[SIZEOF_ENTRY]byte)(unsafe.Pointer(this)))[:], b)
	if this.Proto() != PROTO_ENTRY {
		return nil, ERR_UNKNOWN_PROTO
	}

	this.data = make([]byte, this.size)
	copy(this.data, b[SIZEOF_ENTRY:])
	return this, err
}

func NewEntry() EntryI {
	return &Entry_0_1{}
}
