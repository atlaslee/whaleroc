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
	"time"
)

// -----------------------------------------------------------------------------

const (
	PRTSPVS_CMD_CREATEPARENT = 1
)

type ParentSupervisor struct {
	zsm.Worker
	server  *Server
	backups map[*net.TCPAddr]*Parent
	parent  *Parent
}

func (this *ParentSupervisor) createParent(raddr *net.TCPAddr) *Parent {
	if this.server.Context().BindingAddress().String() == raddr.String() {
		return nil
	}

	if this.parent != nil && this.parent.RemoteAddr().String() == raddr.String() {
		return nil
	}

	for k, _ := range this.backups {
		if k.String() == raddr.String() {
			return nil
		}
	}

	return ParentNew(this.server)
}

func (this *ParentSupervisor) CreateParent(raddr *net.TCPAddr) {
	this.SendMsg3(PRTSPVS_CMD_CREATEPARENT, this, raddr)
}

func (this *ParentSupervisor) PreLoop() (err error) {
	zlog.Debugln("PRTSPVS:", this, "is starting")
	for _, raddr := range this.server.Context().DefaultNodes() {
		if raddr == nil {
			continue
		}
		parent := this.createParent(raddr)
		zsm.WaitForStartupTimeout(parent, 1*time.Second)

		if parent.State() == zsm.STA_RUNNING {
			// v0.1～v0.2只找到1个parent即可
			this.backups[raddr] = parent
			this.parent = parent
			return
		}
	}

	for _, raddr := range this.server.Context().StartupNodes() {
		parent := this.createParent(raddr)
		zsm.WaitForStartupTimeout(parent, 1*time.Second)

		if parent.State() == zsm.STA_RUNNING {
			// v0.1～v0.2只找到1个parent即可
			this.backups[raddr] = parent
			this.parent = parent
			return
		}
	}
	return
}

func (this *ParentSupervisor) AfterLoop() {
	zlog.Debugln("PRTSPVS:", this, "is shutting down")
	for raddr, parent := range this.backups {
		parent.Shutdown()
		delete(this.backups, raddr)
	}
	this.parent = nil
}

func (this *ParentSupervisor) CommandHandle(msg *zsm.Message) (bool, error) {
	return true, nil
}

func ParentSupervisorNew(server *Server) (prtspvs *ParentSupervisor) {
	prtspvs = &ParentSupervisor{
		server:  server,
		backups: server.backups,
		parent:  server.parent}
	prtspvs.Init(prtspvs)
	return
}
