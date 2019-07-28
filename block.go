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
	"crypto/sha256"
	"encoding/base64"
	"time"
	"unsafe"
)

const (
	PROTO_BLOCK        = "BLK"
	SIZEOF_BLOCK_PROTO = len(PROTO_BLOCK)
	SIZEOF_BLOCK       = SIZEOF_BLOCK_PROTO + SIZEOF_VERSION + SIZEOF_HASH + 8 + SIZEOF_HASH + 8 + 1 + SIZEOF_SIGNATURE + 2
)

var (
	PROTO_BLOCK_BYTES = [SIZEOF_BLOCK_PROTO]byte{66, 76, 75}
)

// 区块接口
// 该接口不得修改，以便向前兼容
//
// 用于打包条目以便于传播
type BlockI interface {
	Proto() string     // 用于反序列化校验
	Version() *Version // 用于指定处理的版本
	Id() []byte
	IdString() string
	Height() uint64
	GodHash() [SIZEOF_HASH]byte
	GodHashString() string
	Timestamp() uint64
	EncryptType() byte
	Signature() []byte
	SignatureString() string
	Sign(*Account)
	NumOfEntries() uint16
	Entries() []EntryI
	Bytes() []byte
	SetBytes([]byte) (BlockI, error)
}

// 区块
// 0_1为版本号。发布过的结构不得改变,只能发新结构
// 发布过的属性不得修改，只能新增
type Block_0_1 struct {
	proto        [SIZEOF_BLOCK_PROTO]byte
	version      Version
	id           [SIZEOF_HASH]byte
	height       uint64
	godHash      [SIZEOF_HASH]byte
	timestamp    uint64
	encryptType  byte
	signature    [SIZEOF_SIGNATURE]byte
	numOfEntries uint16
	entries      []EntryI
}

func (this *Block_0_1) Proto() string {
	return string(this.proto[:])
}

func (this *Block_0_1) Version() *Version {
	return &this.version
}

func (this *Block_0_1) Id() []byte {
	return this.id[:]
}

func (this *Block_0_1) IdString() string {
	return base64.URLEncoding.EncodeToString(this.id[:])
}

func (this *Block_0_1) Height() uint64 {
	return this.height
}

func (this *Block_0_1) GodHash() [SIZEOF_HASH]byte {
	return this.godHash
}

func (this *Block_0_1) GodHashString() string {
	return base64.URLEncoding.EncodeToString(this.godHash[:])
}

func (this *Block_0_1) Timestamp() uint64 {
	return this.timestamp
}

func (this *Block_0_1) EncryptType() byte {
	return this.encryptType
}

func (this *Block_0_1) Signature() []byte {
	return this.signature[:]
}

func (this *Block_0_1) Sign(root *Account) {
	copy(this.signature[:], root.Sign(this.Bytes()))
}

func (this *Block_0_1) SignatureString() string {
	return base64.URLEncoding.EncodeToString(this.signature[:])
}

func (this *Block_0_1) NumOfEntries() uint16 {
	return this.numOfEntries
}

func (this *Block_0_1) Entries() []EntryI {
	return this.entries
}

func (this *Block_0_1) Bytes() (b []byte) {
	b = (*(*[SIZEOF_BLOCK]byte)(unsafe.Pointer(this)))[:]

	for _, entry := range this.entries {
		b = append(b, entry.Bytes()...)
	}
	return
}

func (this *Block_0_1) SetBytes(bytes []byte) (block BlockI, err error) {
	n := copy((*(*[SIZEOF_BLOCK]byte)(unsafe.Pointer(this)))[:], bytes)
	if n != SIZEOF_BLOCK {
		return nil, ERR_BUFFER_BROKEN
	}

	if this.Proto() != PROTO_BLOCK {
		return nil, ERR_UNKNOWN_PROTO
	}

	this.entries = make([]EntryI, this.numOfEntries)
	bytes = bytes[SIZEOF_BLOCK:]
	for i := 0; i < int(this.numOfEntries); i++ {
		entry, err := NewEntry().SetBytes(bytes)
		if err != nil {
			return nil, err
		}

		this.entries = append(this.entries, entry)
		bytes = bytes[len(entry.Bytes()):]
	}
	return this, err
}

func GenBlockId(last BlockI) (res [SIZEOF_HASH]byte) {
	data := make([]byte, SIZEOF_HASH+SIZEOF_SIGNATURE)
	copy(data, last.Id())
	copy(data[SIZEOF_HASH:], last.Signature())

	hash := sha256.New()
	hash.Reset()
	hash.Write(data)

	data = hash.Sum(nil)
	copy(res[:], data)
	return
}

func GenGodBlock(god *Account) (block BlockI) {
	block = NewBlock()
	block.Sign(god)
	return
}

func NewBlock() (block BlockI) {
	block = &Block_0_1{
		proto:     PROTO_BLOCK_BYTES,
		version:   *DMT_VERSION,
		timestamp: uint64(time.Now().UnixNano())}
	return
}

func NewBlock1(last BlockI) (block BlockI) {
	block = &Block_0_1{
		id:        GenBlockId(last),
		proto:     PROTO_BLOCK_BYTES,
		version:   *DMT_VERSION,
		timestamp: uint64(time.Now().UnixNano())}
	return
}
