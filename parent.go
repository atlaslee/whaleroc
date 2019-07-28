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

type Parent struct {
	zsm.Monitor
	server     *Server
	conn       *net.TCPConn
	remoteAddr *net.TCPAddr
}

func (this *Parent) RemoteAddr() *net.TCPAddr {
	return this.remoteAddr
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