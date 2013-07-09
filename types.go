/*
Copyright © 2013 mortdeus <mortdeus@gocos2d.org>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
“Software”), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package egl

/*
#include <stdlib.h>
#include <EGL/egl.h>
#include <EGL/eglplatform.h>
*/
import "C"
import (
	"unsafe"
)

type (
	Enum              uint32
	Config            uintptr
	Context           uintptr
	Display           uintptr
	Surface           uintptr
	ClientBuffer      uintptr
	NativeDisplayType unsafe.Pointer
	NativeWindowType  unsafe.Pointer
	NativePixmapType  unsafe.Pointer
)

func goBoolean(n C.EGLBoolean) bool {
	return n == 1
}
func eglBoolean(n bool) C.EGLBoolean {
	var b int
	if n == true {
		b = 1
	}
	return C.EGLBoolean(b)
}

/*
func ProcAdress(proc string) uintptr {

}
*/
