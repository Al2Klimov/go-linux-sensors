package go_linux_sensors

/*
#include "bridge.h"
#cgo LDFLAGS: -lsensors
*/
import "C"

import (
	"os"
	"reflect"
	"sync"
	"unsafe"
)

var errors = map[C.int]string{}
var errorsLock sync.RWMutex

var libsensorsVersion = cStr2str(C.libsensors_version)

var parseError func(err string, lineno int) = nil
var parseErrorWfn func(err, filename string, lineno int) = nil
var fatalError func(proc, err string) = nil

func cStr2str(cStr *C.TCchar) string {
	slen := int(C.strlen(cStr))

	bridge := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cStr)),
		Len:  slen,
		Cap:  slen,
	}

	return string(*(*[]byte)(unsafe.Pointer(&bridge)))
}

func str2cStr(str string) (cStr *C.TCchar, keepAlive []byte) {
	keepAlive = append([]byte(str), 0)
	cStr = (*C.TCchar)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&keepAlive)).Data))
	return
}

//export parseErrorWrapper
func parseErrorWrapper(err *C.TCchar, lineno C.int) {
	parseError(cStr2str(err), int(lineno))
}

//export parseErrorWfnWrapper
func parseErrorWfnWrapper(err, filename *C.TCchar, lineno C.int) {
	parseErrorWfn(cStr2str(err), cStr2str(filename), int(lineno))
}

//export fatalErrorWrapper
func fatalErrorWrapper(proc, err *C.TCchar) {
	fatalError(cStr2str(proc), cStr2str(err))
	os.Exit(1)
}
