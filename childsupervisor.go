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

const (
	CLDSPVS_CMD_CREATECHILD = 1
)

type ChildSupervisor struct {
	zsm.Worker
	server   *Server
	children map[*net.TCPConn]*Child
}

func (this *ChildSupervisor) createChild(conn *net.TCPConn) {
	child := NewChild(this.server, conn)
	child.Startup()
	this.children[conn] = child
}

func (this *ChildSupervisor) CreateChild(conn *net.TCPConn) {
	this.SendMsg3(CLDSPVS_CMD_CREATECHILD, this, conn)
}

func (this *ChildSupervisor) PreLoop() (err error) {
	zlog.Debugln("CLDSPVS:", this, "is starting up")
	return
}

func (this *ChildSupervisor) AfterLoop() {
	zlog.Debugln("CLDSPVS:", this, "is shutting down")
	for conn, child := range this.children {
		child.Shutdown()
		conn.Close()
	}
}

func (this *ChildSupervisor) CommandHandle(msg *zsm.Message) (bool, error) {
	switch msg.Type {
	case CLDSPVS_CMD_CREATECHILD:
		conn, ok := msg.Data.(*net.TCPConn)
		if ok {
			this.createChild(conn)
		}
	}
	return true, nil
}

func NewChildSupervisor(server *Server) (cldspvs *ChildSupervisor) {
	cldspvs = &ChildSupervisor{
		server: server}

	cldspvs.Init(cldspvs)
	return
}
