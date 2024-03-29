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
	PROTO_NODEINFO        = "WRCP/NDIF"
	VER_NODEINFO_STR      = "v0.1.0.0"
	SIZEOF_NODEINFO_PROTO = len(PROTO_NODEINFO)
	SIZEOF_NODEINFO       = SIZEOF_NODEINFO_PROTO + SIZEOF_VERSION + SIZEOF_ADDRESS + SIZEOF_HASH + 8*3
)

var (
	VER_NODEINFO         = NewVersion().SetString(VER_NODEINFO_STR)
	PROTO_NODEINFO_BYTES = [SIZEOF_NODEINFO_PROTO]byte{78, 68, 73, 70}
)

// 节点信息接口
// 该接口不得修改，以便向前兼容
//
// 用于握手，通讯双方判断连接是否成立
type NodeInfoI interface {
	Proto() string     // 用于反序列化校验
	Version() *Version // 用于指定处理的版本
	Address() string
	GodHash() string
	StartupTime() uint64
	CurrentHeight() uint64
	LastUpdateTime() uint64
	Bytes() []byte
	SetBytes([]byte) (NodeInfoI, error)
}

// 节点信息
// 0_1为版本号。发布过的结构不得改变,只能发新结构
// 发布过的属性不得修改，只能新增
type NodeInfo_0_1 struct {
	proto          [SIZEOF_NODEINFO_PROTO]byte
	version        Version
	address        [SIZEOF_ADDRESS]byte
	godHash        [SIZEOF_HASH]byte
	startupTime    uint64
	currentHeight  uint64
	lastUpdateTime uint64
}

func (this *NodeInfo_0_1) Proto() string {
	return string(this.proto[:])
}

func (this *NodeInfo_0_1) Version() *Version {
	return &this.version
}

func (this *NodeInfo_0_1) Address() string {
	return string(this.address[:])
}

func (this *NodeInfo_0_1) GodHash() string {
	return string(this.godHash[:])
}

func (this *NodeInfo_0_1) StartupTime() uint64 {
	return this.startupTime
}

func (this *NodeInfo_0_1) CurrentHeight() uint64 {
	return this.currentHeight
}

func (this *NodeInfo_0_1) LastUpdateTime() uint64 {
	return this.lastUpdateTime
}

func (this *NodeInfo_0_1) Bytes() (b []byte) {
	return (*(*[SIZEOF_NODEINFO]byte)(unsafe.Pointer(this)))[:]
}

func (this *NodeInfo_0_1) SetBytes(bytes []byte) (info NodeInfoI, err error) {
	n := copy((*(*[SIZEOF_NODEINFO]byte)(unsafe.Pointer(this)))[:], bytes)
	if n != SIZEOF_NODEINFO {
		return nil, ERR_BUFFER_BROKEN
	}

	if this.Proto() != PROTO_NODEINFO {
		return nil, ERR_UNKNOWN_PROTO
	}
	return this, err
}

func NewNodeInfo() (info NodeInfoI) {
	return &NodeInfo_0_1{}
}

func NewNodeInfo1(s *Server) (info NodeInfoI) {
	return &NodeInfo_0_1{
		proto:          PROTO_NODEINFO_BYTES,
		version:        *DMT_VERSION,
		address:        s.Id(),
		godHash:        s.Runtime().GodHash,
		startupTime:    s.Runtime().StartupTime,
		currentHeight:  s.Runtime().Height,
		lastUpdateTime: s.Runtime().LastUpdateTime}
}
