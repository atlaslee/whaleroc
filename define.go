/*
MIT License

Copyright (c) 2019 Atlas Lee, 4859345@qq.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package dmt

import (
	"errors"
)

const ( // 版本号
	DMT_VER = "v0.1.0.0"
)

const ( // 算法参数
	BITCOIN_BASE58_TABLE = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

const ( // SIZE参数
	SIZEOF_HASH       = 32
	SIZEOF_BIGINT     = 32
	SIZEOF_RANDOM     = 2 * SIZEOF_BIGINT
	SIZEOF_ADDRESS    = 20
	SIZEOF_PRIVATEKEY = 32
	SIZEOF_PUBLICKEY  = 2 * SIZEOF_BIGINT
	SIZEOF_SIGNATURE  = SIZEOF_PUBLICKEY + 2*SIZEOF_BIGINT
)

const (
	ENCRYPTTYPE_ECDSA = iota
)

var ( // 运行时参数
	DEFAULT_STARTUP_NODES = []string{
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
)

var (
	ERR_UNKNOWN_PROTO   = errors.New("Unknown proto founded")
	ERR_VERSION_TOO_OLD = errors.New("Current version is too old")
)
