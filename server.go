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

const (
	NUMBEROF_BACKUP   = 5
	NUMBEROF_CHILDREN = 10
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

func (this *Server) runUnavailableNodes() {
	for address, lastTime := range this.Context.UnavailableNodes {
		zlog.Traceln(time.Now().UnixNano()-lastTime.UnixNano(), 12*60*60*1000*1000)
		if time.Now().UnixNano()-lastTime.UnixNano() > 12*60*60*1000*1000 {
			node := NodeNew(address, this)
			go node.Run()

			delete(this.Context.UnavailableNodes, address)
		}
	}
}

func (this *Server) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")

	zlog.Traceln("Binding address:", this.Context.BindingAddress.String())
	this.Context.Listener, err = net.ListenTCP("tcp", this.Context.BindingAddress)
	if err != nil {
		return
	}

	this.runUnavailableNodes()

	this.Context.LifeCycle = LIFECYCLE_PARENT_SEARCHING
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

func (this *Server) searchParent() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			return
		default:
			if len(this.Context.Nodes) == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			this.Context.Status = SERVERSTATUS_TOPNODE
			return
		}
	}
}

func (this *Server) Loop() bool {
	if this.Context.LifeCycle == LIFECYCLE_PARENT_SEARCHING {
		zlog.Infoln("Looking for parent nodes.")
		this.searchParent()

		this.Context.LifeCycle = LIFECYCLE_RUNNING
	}
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
	}
	return true
}

func (this *Server) AfterLoop() {
	zlog.Debugln(this, "Shut down.")
}

func (this *Server) nodeStartUpOk(from interface{}) (ok bool) {
	node, ok := from.(*Node)
	if !ok {
		return true
	}

	node.ReceiveState()
	this.Context.Nodes[node.RemoteAddress] = node
	this.Context.ReserveNodes[node.RemoteAddress] = 0
	return
}

func (this *Server) nodeStartUpFailed(from interface{}) (ok bool) {
	node, ok := from.(*Node)
	if !ok {
		return true
	}

	node.ReceiveState()
	this.Context.UnavailableNodes[node.RemoteAddress] = time.Now()
	return
}

func (this *Server) CommandHandle(command int, from, data interface{}) (ok bool) {
	zlog.Traceln("Command", command, from, data, "received.")

	switch command {
	case SERVERCMD_NODE_STARTUP_OK:
		this.nodeStartUpOk(from)

	case SERVERCMD_NODE_STARTUP_FAILED:
		this.nodeStartUpFailed(from)
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
