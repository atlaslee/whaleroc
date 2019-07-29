/*
The MIT License (MIT)

Copyright © 2019 Atlas Lee, 4859345@qq.com

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the “Software”),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
*/

package dmt

import (
	"github.com/atlaslee/zlog"
	"github.com/atlaslee/zsm"
	"net"
)

// -----------------------------------------------------------------------------

// 子节点
// 负责广播数据给子节点
// 同时也接收子节点请求
type Child struct {
	zsm.Worker
	server     *Server
	conn       *net.TCPConn
	remoteAddr *net.TCPAddr
	info       NodeInfoI
}

func (this *Child) PreLoop() (err error) {
	zlog.Debugln("CLD:", this, "is starting up")

	// handshake1: send serverInfo
	err = WriteBytes(this.conn, PROTO_NODEINFO, VER_NODEINFO, NewNodeInfo1(this.server).Bytes())
	if err != nil {
		return
	}

	// handshake2: receive clientinfo
	bytes, err := ReadBytes(this.conn, PROTO_NODEINFO)
	if err != nil {
		return
	}

	this.info = NewNodeInfo1(this.server)
	this.info.SetBytes(bytes)
	return
}

func (this *Child) AfterLoop() {
	zlog.Debugln("CLD:", this, "is shutting down")
	this.conn.Close()
}

func (this *Child) CommandHandle(msg *zsm.Message) (bool, error) {
	return true, nil
}

func NewChild(server *Server, conn *net.TCPConn) (cld *Child) {
	remoteAddr, _ := conn.RemoteAddr().(*net.TCPAddr)
	cld = &Child{
		server:     server,
		conn:       conn,
		remoteAddr: remoteAddr}

	cld.Init(cld)
	return
}
