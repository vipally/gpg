///////////////////////////////////////////////////////////////////
//
// !!!!!!!!!!!! NEVER MODIFY THIS FILE MANUALLY !!!!!!!!!!!!
//
// This file was auto-generated by tool [github.com/vipally/gogp]
// Last update at: [Thu Oct 06 2016 09:12:04]
// Generate from:
//   [github.com/vipally/gogp/examples/example.gp]
//   [github.com/vipally/gogp/examples/example2/example.gpg] [uint]
//
// Tool [github.com/vipally/gogp] info:
// CopyRight 2016 @Ally Dale. All rights reserved.
// Author  : Ally Dale(vipally@gmail.com)
// Blog    : http://blog.csdn.net/vipally
// Site    : https://github.com/vipally
// BuildAt : [Oct  5 2016 22:08:02]
// Version : 2.9.0
// 
///////////////////////////////////////////////////////////////////

//This is an example of using gopg tool for generic-programming
//this is an example of using gopg to define an auto-lock global value with generic type
//it will be realized to real go code by gopg tool through the .gpg file with the same name

package example2

import (
	"sync"
)

//auto locked global value
type AutoLockGblUint struct {
	val  uint
	lock sync.RWMutex
}

//new and init a global value
func NewUint(val uint) *AutoLockGblUint{
	p := &AutoLockGblUint{}
	p.val = val
	return p
}

//get value, if modify is disable, lock is unneeded
func (me *AutoLockGblUint) Get() (r uint) {
	me.lock.RLock()
	defer me.lock.RUnlock()
	r = me.val
	return
}

//set value, if modify is disable, delete this function
func (me *AutoLockGblUint) Set(val uint) (r uint) {
	me.lock.Lock()
	defer me.lock.Unlock()
	r = me.val
	me.val = val
	return
}
