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

const (
	SERVERCMD_NODE_STARTUP_OK     = 100
	SERVERCMD_NODE_STARTUP_FAILED = 101
)

var SERVERSTATUSES []string = []string{
	"SERVERSTATUS_ROOT",
	"SERVERSTATUS_TOPNODE",
	"SERVERSTATUS_NODE",
	"SERVERSTATUS_LEAF",
	"SERVERSTATUS_CLIENT"}

type Server struct {
	zsm.StateMachine
	Context *Context
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

func (this *Server) waitForParent() {
}

func (this *Server) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")

	zlog.Traceln("Binding address:", this.Context.BindingAddress.String())
	this.Context.Listener, err = net.ListenTCP("tcp", this.Context.BindingAddress)
	if err != nil {
		return
	}

	this.runUnusedReserveNodes()

	this.waitForParent()
	return
}

func (this *Server) rootLoop() bool {
	zlog.Traceln("Loop as root.")
	return true
}

func (this *Server) topNodeLoop() bool {
	zlog.Traceln("Loop as top node.")
	return true
}

func (this *Server) nodeLoop() bool {
	zlog.Traceln("Loop as node.")
	return true
}

func (this *Server) leafLoop() bool {
	zlog.Traceln("Loop as leaf.")
	return true
}

func (this *Server) clientLoop() bool {
	zlog.Traceln("Loop as client.")
	return true
}

func (this *Server) Loop() bool {
	zlog.Traceln("Looping.")
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
	zlog.Debugln(this, "Shut down.")
}

func (this *Server) nodeStartUpOK(node *Node) (ok bool) {
	node.ReceiveState()
	return true
}

func (this *Server) nodeStartUpFailed(node *Node) (ok bool) {
	//node.ReceiveState()
	return true
}

func (this *Server) CommandHandle(command int, from, data interface{}) (ok bool) {
	zlog.Traceln("Command", command, from, data, "received.")

	var node *Node = nil
	if from != nil {
		node, _ = from.(*Node)
	}

	switch command {
	case SERVERCMD_NODE_STARTUP_OK:
		ok = this.nodeStartUpOK(node)
	case SERVERCMD_NODE_STARTUP_FAILED:
		ok = this.nodeStartUpFailed(node)
	default:
	}
	return ok
}

func ServerNew(context *Context) (server *Server) {
	server = &Server{
		StateMachine: zsm.StateMachine{},
		Context:      context}

	server.Init(server)
	return
}
