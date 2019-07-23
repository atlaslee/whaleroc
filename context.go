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
	"net"
	"time"
)

const (
	LIFECYCLE_INITIALIZING = iota
	LIFECYCLE_PARENT_SEARCHING
	LIFECYCLE_RUNNING
	LIFECYCLE_SHUTTING
)

var LIFECYCLES []string = []string{
	"LIFECYCLE_INITIALIZING",
	"LIFECYCLE_PARENT_SEARCHING",
	"LIFECYCLE_RUNNING",
	"LIFECYCLE_SHUTTING"}

const (
	SERVERSTATUS_ROOT = iota
	SERVERSTATUS_TOPNODE
	SERVERSTATUS_NODE
	SERVERSTATUS_LEAF
	SERVERSTATUS_CLIENT
)

type Context struct {
	BindingAddress   *net.TCPAddr
	Blocks           *list.List
	LifeCycle        uint8
	Listener         *net.TCPListener
	Nodes            map[*net.TCPAddr]*Node
	RecommandNodes   *list.List
	ReserveNodes     map[*net.TCPAddr]uint64
	Status           uint8
	UnavailableNodes map[*net.TCPAddr]time.Time
}

func ContextNew(bindingAddress string, reserveNodes []string) (context *Context) {
	tcpAddress, err := net.ResolveTCPAddr("tcp", bindingAddress)
	if err != nil {
		return
	}

	context = &Context{
		BindingAddress:   tcpAddress,
		Blocks:           list.New(),
		LifeCycle:        LIFECYCLE_INITIALIZING,
		Nodes:            make(map[*net.TCPAddr]*Node),
		RecommandNodes:   list.New(),
		ReserveNodes:     make(map[*net.TCPAddr]uint64),
		Status:           SERVERSTATUS_ROOT,
		UnavailableNodes: make(map[*net.TCPAddr]time.Time)}

	for n := 0; n < len(reserveNodes); n++ {
		if bindingAddress == reserveNodes[n] {
			zlog.Warningln("Cannot connect to self. Skipped.")
			continue
		}

		address, err := net.ResolveTCPAddr("tcp", reserveNodes[n])
		if err != nil {
			zlog.Warningln("Failed in resolving", reserveNodes[n], "to tcp address. Skipped.")
			continue
		}

		context.UnavailableNodes[address] = time.Unix(0, 0)
	}
	return
}
