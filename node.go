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
	NODECMD_SETROOT    = 200
	NODECMD_SETTOPNODE = 201
	NODECMD_SETNODE    = 202
	NODECMD_SETLEAF    = 203
)

func (this *Node) backupRootLoop() bool {
	zlog.Traceln("Backup as root.")
	return true
}

func (this *Node) backupTopNodeLoop() bool {
	zlog.Traceln("Backup as top node.")
	return true
}

func (this *Node) backupNodeLoop() bool {
	zlog.Traceln("Backup as node.")
	return true
}

func (this *Node) backupLeafLoop() bool {
	zlog.Traceln("Backup as leaf.")
	return true
}

func (this *Node) backupClientLoop() bool {
	zlog.Traceln("Backup as client.")
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
	zlog.Traceln("Parent as root.")
	return true
}

func (this *Node) parentTopNodeLoop() bool {
	zlog.Traceln("Parent as top node.")
	return true
}

func (this *Node) parentNodeLoop() bool {
	zlog.Traceln("Parent as node.")
	return true
}

func (this *Node) parentLeafLoop() bool {
	zlog.Traceln("Parent as leaf.")
	return true
}

func (this *Node) parentClientLoop() bool {
	zlog.Traceln("Parent as client.")
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
	zlog.Traceln("Child as root.")
	return true
}

func (this *Node) childTopNodeLoop() bool {
	zlog.Traceln("Child as root.")
	return true
}

func (this *Node) childNodeLoop() bool {
	zlog.Traceln("Child as root.")
	return true
}

func (this *Node) childLeafLoop() bool {
	zlog.Traceln("Child as root.")
	return true
}

func (this *Node) childClientLoop() bool {
	zlog.Traceln("Child as root.")
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
	zlog.Debugln(this, "Starting up.")

	zlog.Traceln("Connecting to address:", this.RemoteAddress.String())
	this.Connection, err = net.DialTCP("tcp", nil, this.RemoteAddress)
	if err != nil {
		zlog.Warningln(this.RemoteAddress.String(), "connect failed.")
		this.Server.SendCommand2(SERVERCMD_NODE_STARTUP_FAILED, this)
		return
	}

	zlog.Infoln("Connection to", this.RemoteAddress.String(), "established.")
	this.LifeCycle = LIFECYCLE_RUNNING
	this.Server.SendCommand2(SERVERCMD_NODE_STARTUP_OK, this)
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
	zlog.Debugln(this, "Shut down.")
}

func (this *Node) setRoot() bool {
	zlog.Infoln("Run as root.")
	return true
}

func (this *Node) setTopNode() bool {
	zlog.Infoln("Run as top node.")
	return true
}

func (this *Node) setNode() bool {
	zlog.Infoln("Run as node.")
	return true
}

func (this *Node) setLeaf() bool {
	zlog.Infoln("Run as leaf.")
	return true
}

func (this *Node) CommandHandle(command int, from, data interface{}) (ok bool) {
	zlog.Traceln("Command", command, from, data, "received.")

	//node, _ := value.(*Node)

	switch command {
	case NODECMD_SETROOT:
		ok = this.setRoot()
	case NODECMD_SETTOPNODE:
		ok = this.setTopNode()
	case NODECMD_SETNODE:
		ok = this.setNode()
	case NODECMD_SETLEAF:
		ok = this.setLeaf()
	default:
		return true
	}
	return
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
