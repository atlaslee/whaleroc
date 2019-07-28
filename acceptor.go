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
)

// -----------------------------------------------------------------------------

// 接收器
// 负责监听端口，接收外部请求并通知子节点监管者
type Acceptor struct {
	zsm.Monitor
	server *Server
}

func (this *Acceptor) PreLoop() (err error) {
	zlog.Debugln(this, "Starting up.")
	return
}

func (this *Acceptor) Loop() (ok bool, err error) {
	// 监听Listener
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
	this.server.Listener().Close()
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
