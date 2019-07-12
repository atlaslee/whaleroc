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
	LIFECYCLE_INITIALIZING = iota
	LIFECYCLE_RUNNING
	LIFECYCLE_PAUSING
	LIFECYCLE_RESUMING
	LIFECYCLE_SHUTTING
)

var LIFECYCLES []string = []string{
	"LIFECYCLE_INITIALIZING",
	"LIFECYCLE_RUNNING",
	"LIFECYCLE_PAUSING",
	"LIFECYCLE_RESUMING",
	"LIFECYCLE_SHUTTING"}

const (
	NODESTATUS_BACKUP = iota
	NODESTATUS_PARENT
	NODESTATUS_CHILD
)

var NODESTATUSES []string = []string{
	"NODESTATUS_BACKUP",
	"NODESTATUS_PARENT",
	"NODESTATUS_CHILD"}

type Node struct {
	zsm.StateMachine
	Blocks        *list.List
	Connection    *net.TCPConn
	LifeCycle     uint8
	RemoteAddress *net.TCPAddr
	Server        *Server
	Status        uint8
}

const (
	NODE_STARTUP_OK     = 100
	NODE_STARTUP_FAILED = 101
)

func (this *Node) backupRootLoop() bool {
	return true
}

func (this *Node) backupTopNodeLoop() bool {
	return true
}

func (this *Node) backupNodeLoop() bool {
	return true
}

func (this *Node) backupLeafLoop() bool {
	return true
}

func (this *Node) backupClientLoop() bool {
	return true
}

func (this *Node) backupLoop() bool {
	switch this.Server.Context.Status {
	case SERVERSTATUS_ROOT:
		return this.backupRootLoop()
	case SERVERSTATUS_TOPNODE:
		return this.backupTopNodeLoop()
	case SERVERSTATUS_NODE:
		return this.backupNodeLoop()
	case SERVERSTATUS_LEAF:
		return this.backupLeafLoop()
	case SERVERSTATUS_CLIENT:
		return this.backupClientLoop()
	default:
		return true
	}
}

func (this *Node) parentRootLoop() bool {
	return true
}

func (this *Node) parentTopNodeLoop() bool {
	return true
}

func (this *Node) parentNodeLoop() bool {
	return true
}

func (this *Node) parentLeafLoop() bool {
	return true
}

func (this *Node) parentClientLoop() bool {
	return true
}

func (this *Node) parentLoop() bool {
	switch this.Server.Context.Status {
	case SERVERSTATUS_ROOT:
		return this.parentRootLoop()
	case SERVERSTATUS_TOPNODE:
		return this.parentTopNodeLoop()
	case SERVERSTATUS_NODE:
		return this.parentNodeLoop()
	case SERVERSTATUS_LEAF:
		return this.parentLeafLoop()
	case SERVERSTATUS_CLIENT:
		return this.parentClientLoop()
	default:
		return true
	}
}

func (this *Node) childRootLoop() bool {
	return true
}

func (this *Node) childTopNodeLoop() bool {
	return true
}

func (this *Node) childNodeLoop() bool {
	return true
}

func (this *Node) childLeafLoop() bool {
	return true
}

func (this *Node) childClientLoop() bool {
	return true
}

func (this *Node) childLoop() bool {
	switch this.Server.Context.Status {
	case SERVERSTATUS_ROOT:
		return this.childRootLoop()
	case SERVERSTATUS_TOPNODE:
		return this.childTopNodeLoop()
	case SERVERSTATUS_NODE:
		return this.childNodeLoop()
	case SERVERSTATUS_LEAF:
		return this.childLeafLoop()
	case SERVERSTATUS_CLIENT:
		return this.childClientLoop()
	default:
		return true
	}
}

func (this *Node) PreLoop() (err error) {
	zlog.Debugln("Starting up.")

	this.Connection, err = net.DialTCP("tcp", nil, this.RemoteAddress)
	if err != nil {
		this.Server.SendCommand2(NODE_STARTUP_FAILED, this)
		return
	}

	this.LifeCycle = LIFECYCLE_RUNNING
	this.Server.SendCommand2(NODE_STARTUP_OK, this)
	return
}

func (this *Node) Loop() bool {
	switch this.Status {
	case NODESTATUS_BACKUP:
		return this.backupLoop()
	case NODESTATUS_PARENT:
		return this.parentLoop()
	case NODESTATUS_CHILD:
		return this.childLoop()
	}
	return true
}

func (this *Node) AfterLoop() {
	zlog.Debugln("Shut down.")
}

func (this *Node) CommandHandle(command int, value interface{}) bool {
	zlog.Traceln("Command", command, "received.")
	return true
}

func NodeNew(address *net.TCPAddr, server *Server) (node *Node) {
	node = &Node{
		StateMachine:  zsm.StateMachine{},
		Blocks:        list.New(),
		RemoteAddress: address,
		Server:        server}

	node.Init(node)
	return
}
