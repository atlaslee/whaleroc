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
	"container/list"
	"github.com/atlaslee/zlog"
	"github.com/atlaslee/zsm"
	"net"
)

const (
	NUMBEROF_BACKUP   = 5
	NUMBEROF_CHILDREN = 10
)

const (
	SERVERSTATUS_ROOT = iota
	SERVERSTATUS_TOPNODE
	SERVERSTATUS_NODE
	SERVERSTATUS_LEAF
	SERVERSTATUS_CLIENT
)

var SERVERSTATUSES []string = []string{
	"SERVERSTATUS_ROOT",
	"SERVERSTATUS_TOPNODE",
	"SERVERSTATUS_NODE",
	"SERVERSTATUS_LEAF",
	"SERVERSTATUS_CLIENT"}

type Context struct {
	BindingAddress     *net.TCPAddr
	Blocks             *list.List
	LifeCycle          uint8
	Listener           *net.TCPListener
	Nodes              map[*net.TCPAddr]*Node
	RecommandNodes     *list.List
	ReserveNodes       map[*net.TCPAddr]uint64
	Status             uint8
	UnusedReserveNodes *list.List
}

func ContextNew(bindingAddress *net.TCPAddr, reserveNodes []*net.TCPAddr) (context *Context) {
	context = &Context{
		BindingAddress:     bindingAddress,
		Blocks:             list.New(),
		LifeCycle:          LIFECYCLE_INITIALIZING,
		Nodes:              make(map[*net.TCPAddr]*Node, 0),
		RecommandNodes:     list.New(),
		ReserveNodes:       make(map[*net.TCPAddr]uint64),
		Status:             SERVERSTATUS_ROOT,
		UnusedReserveNodes: list.New()}
	for n := 0; n < len(reserveNodes); n++ {
		context.UnusedReserveNodes.PushBack(reserveNodes[n])
	}
	return
}

type Server struct {
	zsm.StateMachine
	Context *Context
}

func (this *Server) nodeStartUpOK(node *Node) (ok bool) {
	node.ReceiveState()
	return
}

func (this *Server) nodeStartUpFailed(node *Node) (ok bool) {
	node.ReceiveState()
	return
}

func (this *Server) runUnusedReserveNodes() {
	unusedReserveNodeElement := this.Context.UnusedReserveNodes.Front()
	for n := 0; n < this.Context.UnusedReserveNodes.Len(); n++ {
		address, ok := unusedReserveNodeElement.Value.(*net.TCPAddr)
		if ok {
			node := NodeNew(address, this)
			go node.Run()
		}

		unusedReserveNodeElement = unusedReserveNodeElement.Next()
	}
}

func (this *Server) rootLoop() bool {
	return true
}

func (this *Server) topNodeLoop() bool {
	return true
}

func (this *Server) nodeLoop() bool {
	return true
}

func (this *Server) leafLoop() bool {
	return true
}

func (this *Server) clientLoop() bool {
	return true
}

func (this *Server) PreLoop() (err error) {
	zlog.Debugln("Starting up.")

	this.Context.Listener, err = net.ListenTCP("tcp", this.Context.BindingAddress)
	if err != nil {
		return
	}

	this.runUnusedReserveNodes()
	return
}

func (this *Server) Loop() bool {
	switch this.Context.Status {
	case SERVERSTATUS_ROOT:
		return this.rootLoop()
	case SERVERSTATUS_TOPNODE:
		return this.topNodeLoop()
	case SERVERSTATUS_NODE:
		return this.nodeLoop()
	case SERVERSTATUS_LEAF:
		return this.leafLoop()
	case SERVERSTATUS_CLIENT:
		return this.clientLoop()
	default:
		return true
	}
}

func (this *Server) AfterLoop() {
	zlog.Debugln("Shut down.")
}

func (this *Server) CommandHandle(command int, value interface{}) (ok bool) {
	zlog.Traceln("Command", command, "received.")

	node, _ := value.(*Node)

	switch command {
	case NODE_STARTUP_OK:
		ok = this.nodeStartUpOK(node)
	case NODE_STARTUP_FAILED:
		ok = this.nodeStartUpFailed(node)
	default:
	}
	return
}

func ServerNew(context *Context) (server *Server) {
	server = &Server{
		StateMachine: zsm.StateMachine{},
		Context:      context}

	server.Init(server)
	return
}
