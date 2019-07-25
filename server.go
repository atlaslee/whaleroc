/* The MIT License (MIT)
Copyright © 2018 by Atlas Lee(atlas@fpay.io)

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

type Acceptor struct {
	zsm.StateMachine
	server *Server
}

func (this *Acceptor) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")
	return
}

func (this *Acceptor) Loop() (ok bool, err error) {
	conn, err := this.server.Listener().AcceptTCP()
	if err != nil {
		zlog.Fatalln("Acceptor:", this, "failed:", err)
		ok = false
		return
	}
	this.server.ChildSupervisor().CreateChild(conn)
	return
}

func (this *Acceptor) AfterLoop() {
}

func (this *Acceptor) CommandHandle(message *zsm.Message) (bool, error) {
	return true, nil
}

func AcceptorNew(server *Server) (acceptor *Acceptor) {
	acceptor = &Acceptor{
		server: server}

	acceptor.Init(acceptor)
	return
}

// -----------------------------------------------------------------------------

type Child struct {
	zsm.StateMachine
	server     *Server
	conn       *net.TCPConn
	remoteAddr *net.TCPAddr
}

func (this *Child) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")
	return
}

func (this *Child) Loop() (ok bool, err error) {
	<-time.After(100 * time.Millisecond)
	return
}

func (this *Child) AfterLoop() {
}

func (this *Child) CommandHandle(msg *zsm.Message) (bool, error) {
	return true, nil
}

func ChildNew(server *Server) (cld *Child) {
	cld = &Child{
		server: server}

	cld.Init(cld)
	return
}

// -----------------------------------------------------------------------------

const (
	CLDSPVS_CMD_CREATECHILD = 1
)

type ChildSupervisor struct {
	zsm.StateMachine
	server   *Server
	children map[*net.TCPConn]*Child
}

func (this *ChildSupervisor) CreateChild(conn *net.TCPConn) {
	this.SendMessage3(CLDSPVS_CMD_CREATECHILD, this, conn)
}

func (this *ChildSupervisor) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")
	return
}

func (this *ChildSupervisor) Loop() (ok bool, err error) {
	<-time.After(100 * time.Millisecond)
	return
}

func (this *ChildSupervisor) AfterLoop() {
}

func (this *ChildSupervisor) CommandHandle(msg *zsm.Message) (bool, error) {
	return true, nil
}

func ChildSupervisorNew(server *Server) (cldspvs *ChildSupervisor) {
	cldspvs = &ChildSupervisor{
		server: server}

	cldspvs.Init(cldspvs)
	return
}

// -----------------------------------------------------------------------------

var DEFAULT_STARTUP_NODES []string = []string{
	"127.0.0.1:8000",
	"127.0.0.1:8001",
	"127.0.0.1:8002",
	"127.0.0.1:8003",
	"127.0.0.1:8004",
	"127.0.0.1:8005",
	"127.0.0.1:8006",
	"127.0.0.1:8007",
	"127.0.0.1:8008",
	"127.0.0.1:8009"}

type Context struct {
	bindingAddress *net.TCPAddr
	defaultNodes   []*net.TCPAddr
	startupNodes   []*net.TCPAddr
}

func (this *Context) BindingAddress() *net.TCPAddr {
	return this.bindingAddress
}

func (this *Context) DefaultNodes() []*net.TCPAddr {
	return this.defaultNodes
}

func (this *Context) StartupNodes() []*net.TCPAddr {
	return this.startupNodes
}

func ContextNew(bindingAddress string, startupNodes []string) (ctx *Context) {
	tcpAddress, err := net.ResolveTCPAddr("tcp", bindingAddress)
	if err != nil {
		zlog.Fatalln("Ctx:", nil, "failed:", err)
		return
	}

	ctx = &Context{
		bindingAddress: tcpAddress,
		defaultNodes:   make([]*net.TCPAddr, 10),
		startupNodes:   make([]*net.TCPAddr, 10)}

	for n := 0; n < len(DEFAULT_STARTUP_NODES); n++ {
		if DEFAULT_STARTUP_NODES[n] == bindingAddress {
			zlog.Warningln("Ctx:", ctx, "skipped:", DEFAULT_STARTUP_NODES[n])
			continue
		}

		address, err := net.ResolveTCPAddr("tcp", DEFAULT_STARTUP_NODES[n])
		if err != nil {
			zlog.Warningln("Ctx:", ctx, "skipped:", DEFAULT_STARTUP_NODES[n])
			continue
		}

		ctx.defaultNodes = append(ctx.defaultNodes, address)
	}

	for n := 0; n < len(startupNodes); n++ {
		if startupNodes[n] == bindingAddress {
			zlog.Warningln("Ctx:", ctx, "skipped:", startupNodes[n])
			continue
		}

		address, err := net.ResolveTCPAddr("tcp", startupNodes[n])
		if err != nil {
			zlog.Warningln("Ctx:", ctx, "skipped:", startupNodes[n])
			continue
		}

		ctx.startupNodes = append(ctx.startupNodes, address)
	}
	return
}

// -----------------------------------------------------------------------------

type Parent struct {
	zsm.StateMachine
	server     *Server
	conn       *net.TCPConn
	remoteAddr *net.TCPAddr
}

func (this *Parent) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")
	return
}

func (this *Parent) Loop() (ok bool, err error) {
	<-time.After(100 * time.Millisecond)
	return
}

func (this *Parent) AfterLoop() {
}

func (this *Parent) CommandHandle(msg *zsm.Message) (bool, error) {
	return true, nil
}

func ParentNew(server *Server) (prt *Parent) {
	prt = &Parent{
		server: server}

	prt.Init(prt)
	return
}

// -----------------------------------------------------------------------------

type ParentSupervisor struct {
	zsm.StateMachine
	server  *Server
	backups map[*net.TCPAddr]*Parent
	parent  *Parent
}

func (this *ParentSupervisor) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")
	return
}

func (this *ParentSupervisor) Loop() (ok bool, err error) {
	<-time.After(100 * time.Millisecond)
	return
}

func (this *ParentSupervisor) AfterLoop() {
}

func (this *ParentSupervisor) CommandHandle(msg *zsm.Message) (bool, error) {
	return true, nil
}

func ParentSupervisorNew(server *Server) (prtspvs *ParentSupervisor) {
	prtspvs = &ParentSupervisor{
		server: server}

	prtspvs.Init(prtspvs)
	return
}

// -----------------------------------------------------------------------------

type Server struct {
	zsm.StateMachine
	id               int
	context          *Context
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

func (this *Server) ID() int {
	return this.id
}

func (this *Server) Listener() *net.TCPListener {
	return this.listener
}

func (this *Server) Parent() (parent *Parent) {
	return this.parent
}

func (this *Server) PreLoop() (err error) {
	zlog.Debugln("SVR:", this.id, "starting up")

	zlog.Traceln("SVR:", this.id, "binding address:", this.context.BindingAddress().String())
	this.listener, err = net.ListenTCP("tcp", this.context.BindingAddress())
	if err != nil {
		zlog.Fatalln("SVR:", this.id, "failed:", err)
		return
	}

	this.parentSupervisor = ParentSupervisorNew(this)
	this.parentSupervisor.Run()

	zsm.WaitForStartupTimeout(this.parentSupervisor, 5*time.Second)

	this.childSupervisor = ChildSupervisorNew(this)
	this.childSupervisor.Run()
	return
}

func (this *Server) Loop() (bool, error) {
	<-time.After(100 * time.Millisecond)
	return true, nil
}

func (this *Server) AfterLoop() {
	zlog.Debugln("SVR:", this.id, "shutdown")
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
