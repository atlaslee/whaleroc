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
	"encoding/base64"
	"github.com/atlaslee/zlog"
	"github.com/atlaslee/zsm"
	"net"
	"time"
)

// -----------------------------------------------------------------------------

type Server struct {
	zsm.Worker
	address          [SIZEOF_ADDRESS]byte
	context          *Context
	runtime          *Runtime
	listener         *net.TCPListener
	acceptor         *Acceptor
	parent           *Parent
	backups          map[*net.TCPAddr]*Parent
	parentSupervisor *ParentSupervisor
	children         map[*net.TCPConn]*Child
	childSupervisor  *ChildSupervisor
}

func (this *Server) Acceptor() *Acceptor {
	return this.acceptor
}

func (this *Server) Backup(remoteAddress *net.TCPAddr) (backup *Parent) {
	backup, _ = this.backups[remoteAddress]
	return
}

func (this *Server) Backups() map[*net.TCPAddr]*Parent {
	return this.backups
}

func (this *Server) Children() map[*net.TCPConn]*Child {
	return this.children
}

func (this *Server) ChildSupervisor() *ChildSupervisor {
	return this.childSupervisor
}

func (this *Server) Context() *Context {
	return this.context
}

func (this *Server) Runtime() *Runtime {
	return this.runtime
}

func (this *Server) Address() [SIZEOF_ADDRESS]byte {
	return this.address
}

func (this *Server) AddressString() string {
	return base64.URLEncoding.EncodeToString(this.address[:])
}

func (this *Server) Listener() *net.TCPListener {
	return this.listener
}

func (this *Server) Parent() (parent *Parent) {
	return this.parent
}

func (this *Server) PreLoop() (err error) {
	zlog.Debugln("SVR:", this.AddressString(), "starting")

	zlog.Traceln("SVR:", this.AddressString(), "binding address:", this.context.BindingAddress().String())
	this.listener, err = net.ListenTCP("tcp", this.context.BindingAddress())
	if err != nil {
		zlog.Fatalln("SVR:", this.AddressString(), "failed:", err)
		return
	}

	this.parentSupervisor = ParentSupervisorNew(this)
	this.parentSupervisor.Run()

	zsm.WaitForStartupTimeout(this.parentSupervisor, 5*time.Second)

	this.childSupervisor = ChildSupervisorNew(this)
	this.childSupervisor.Run()
	return
}

func (this *Server) AfterLoop() {
	zlog.Debugln("SVR:", this.AddressString(), "stopping")
}

func (this *Server) CommandHandle(message *zsm.Message) (bool, error) {
	return true, nil
}

func ServerNew(context *Context) (server *Server) {
	server = &Server{
		context:  context,
		backups:  make(map[*net.TCPAddr]*Parent),
		children: make(map[*net.TCPConn]*Child)}

	server.Init(server)
	return
}

// -----------------------------------------------------------------------------
