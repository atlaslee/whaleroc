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
	"strconv"
	"strings"
	"unsafe"
)

const (
	SIZEOF_VERSION = 4
)

// 版本
type Version struct {
	main, milestone, minor, build uint8
}

func atou8(str string) uint8 {
	i, err := strconv.Atoi(str)
	if err != nil {
		i = 0
	}
	return uint8(i)
}

func u8toa(u uint8) string {
	return strconv.Itoa(int(u))
}

func (this *Version) Main() int {
	return int(this.main)
}

func (this *Version) SetMain(main int) {
	this.main = uint8(main)
}

func (this *Version) Milestone() int {
	return int(this.milestone)
}

func (this *Version) SetMilestone(milestone int) {
	this.milestone = uint8(milestone)
}

func (this *Version) Minor() int {
	return int(this.minor)
}

func (this *Version) SetMinor(minor int) {
	this.minor = uint8(minor)
}

func (this *Version) Build() int {
	return int(this.build)
}

func (this *Version) SetBuild(build int) {
	this.build = uint8(build)
}

// 与另一个版本比较
// 大于0则本版本更新
// 等于0则两者版本一致
// 小于0则本版本更旧
func (this *Version) Cmp(v *Version) (i int) {
	i = int(this.main) - int(v.main)
	if i != 0 {
		return
	}

	i = int(this.milestone) - int(v.milestone)
	if i != 0 {
		return
	}

	i = int(this.minor) - int(v.minor)
	if i != 0 {
		return
	}

	i = int(this.build) - int(v.build)
	return
}

func (this *Version) Newer(v *Version) bool {
	return this.Cmp(v) > 0
}

func (this *Version) NotNewer(v *Version) bool {
	return this.Cmp(v) <= 0
}

func (this *Version) Older(v *Version) bool {
	return this.Cmp(v) < 0
}

func (this *Version) NotOlder(v *Version) bool {
	return this.Cmp(v) >= 0
}

func (this *Version) Equal(v *Version) bool {
	return this.Cmp(v) == 0
}

func (this *Version) Bytes() []byte {
	return (*[SIZEOF_VERSION]byte)(unsafe.Pointer(this))[:]
}

func (this *Version) SetBytes(bytes []byte) {
	copy(this.Bytes(), bytes)
}

func (this *Version) String() string {
	return "v" + u8toa(this.main) + "." + u8toa(this.milestone) + "." + u8toa(this.minor) + "." + u8toa(this.build)
}

func (this *Version) SetString(str string) *Version {
	strs := strings.Split(str, ".")
	num := len(strs)

	if num > 0 {
		this.main = atou8(strings.Trim(strs[0], "v"))
	} else {
		return this
	}

	if num > 1 {
		this.milestone = atou8(strs[1])
	} else {
		return this
	}

	if num > 2 {
		this.minor = atou8(strs[2])
	} else {
		return this
	}

	if num > 3 {
		this.build = atou8(strs[3])
	}
	return this
}

func NewVersion() *Version {
	return &Version{}
}

func NewVersion1(main int) *Version {
	return &Version{uint8(main), 0, 0, 0}
}

func NewVersion2(main, milestone int) *Version {
	return &Version{uint8(main), uint8(milestone), 0, 0}
}

func NewVersion3(main, milestone, minor int) *Version {
	return &Version{uint8(main), uint8(milestone), uint8(minor), 0}
}

func NewVersion4(main, milestone, minor, build int) *Version {
	return &Version{uint8(main), uint8(milestone), uint8(minor), uint8(build)}
}
