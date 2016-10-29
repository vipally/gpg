///////////////////////////////////////////////////////////////////
//
// !!!!!!!!!!!! NEVER MODIFY THIS FILE MANUALLY !!!!!!!!!!!!
//
// This file was auto-generated by tool [github.com/vipally/gogp]
// Last update at: [Sat Oct 29 2016 17:23:54]
// Generate from:
//   [github.com/vipally/gogp/examples/stack.gp_string.go]
//   [github.com/vipally/gogp/examples/example.gpg] [stack_string]
//
// Tool [github.com/vipally/gogp] info:
// CopyRight 2016 @Ally Dale. All rights reserved.
// Author  : Ally Dale(vipally@gmail.com)
// Blog    : http://blog.csdn.net/vipally
// Site    : https://github.com/vipally
// BuildAt : [Oct 24 2016 20:25:45]
// Version : 3.0.0.final
//
///////////////////////////////////////////////////////////////////

//this file is used to import by other gp files
//it cannot use independently
//simulation C++ stl functors
package examples

type Comparerstring interface {
	F(left, right string) bool
}

type ComparerstringCreator int

const (
	LESSER_string ComparerstringCreator = iota
	GREATER_string
)

func (me ComparerstringCreator) Create() (cmp Comparerstring) {
	switch me {
	case LESSER_string:
		cmp = Lesserstring(0)
	case GREATER_string:
		cmp = Greaterstring(0)
	}
	return
}

type Lesserstring byte

func (this Lesserstring) F(left, right string) (ok bool) {

	ok = left < right

	return
}

type Greaterstring byte

func (this Greaterstring) F(left, right string) (ok bool) {

	ok = left > right

	return
}