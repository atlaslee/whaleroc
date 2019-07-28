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
	"net"
)

// -----------------------------------------------------------------------------

type Context struct {
	version        *Version
	account        *Account
	godBlock       BlockI
	bindingAddress *net.TCPAddr
	defaultNodes   []*net.TCPAddr
	startupNodes   []*net.TCPAddr
}

func (this *Context) Version() *Version {
	return this.version
}

func (this *Context) Account() *Account {
	return this.account
}

func (this *Context) GodBlock() BlockI {
	return this.godBlock
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

func ContextNew(
	bindingAddress string,
	account *Account,
	godBlock BlockI,
	startupNodes []string) (ctx *Context) {

	tcpAddress, err := net.ResolveTCPAddr("tcp", bindingAddress)
	if err != nil {
		zlog.Fatalln("Ctx:", nil, "failed:", err)
		return
	}

	ctx = &Context{
		version:        NewVersion().SetString(DMT_VER),
		account:        account,
		godBlock:       godBlock,
		bindingAddress: tcpAddress,
		defaultNodes:   make([]*net.TCPAddr, 1),
		startupNodes:   make([]*net.TCPAddr, 1)}

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
