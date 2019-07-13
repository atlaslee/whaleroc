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
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	zlog.Traceln("TestServer")

	context := ContextNew("127.0.0.1:8000",
		[]string{
			"127.0.0.1:8000",
			"127.0.0.1:8001",
			/*"127.0.0.1:8002",
			"127.0.0.1:8003",
			"127.0.0.1:8004",
			"127.0.0.1:8005",
			"127.0.0.1:8006",
			"127.0.0.1:8007",
			"127.0.0.1:8008",*/
			"127.0.0.1:8009"})

	server := ServerNew(context)
	server.Startup()
	time.Sleep(2 * time.Second)
	server.Shutdown()
}
