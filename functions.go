package go_linux_sensors

/*
#include "bridge.h"
#cgo LDFLAGS: -lsensors
*/
import "C"

import (
	"os"
	"runtime"
	"unsafe"
)

func Init(input *os.File) error {
	file := (*C.Tfile)(nil)

	if input != nil {
		fd, errDp := C.dup(C.int(input.Fd()))
		if fd == C.int(-1) {
			return errDp
		}

		roMode := [2]byte{'r', 0}
		var errFO error

		file, errFO = C.fdopen(fd, (*C.TCchar)(unsafe.Pointer(&roMode)))
		if file == nil {
			C.close(fd)
			return errFO
		}

		defer C.fclose(file)
	}

	if errSI := C.sensors_init(file); errSI != C.int(0) {
		return Error{errSI}
	}

	return nil
}

func Cleanup() {
	C.sensors_cleanup()
}

func GetLibsensorsVersion() string {
	return libsensorsVersion
}

func ParseChipName(origName string) (*ChipName, error) {
	origNameCStr, origNameCStrKeepAlive := str2cStr(origName)
	res := &ChipName{}
	res.freeOnClose = true

	defer runtime.KeepAlive(origNameCStrKeepAlive)

	if errPCN := C.sensors_parse_chip_name(origNameCStr, &res.chipName); errPCN != C.int(0) {
		return nil, Error{errPCN}
	}

	return res, nil
}

func GetDetectedChips(match *ChipName) []*ChipName {
	actualMatch := (*C.struct_sensors_chip_name)(nil)
	nr := C.int(0)

	if match != nil {
		actualMatch = &match.chipName
	}

	chips := []*ChipName{}

	for {
		chip := C.sensors_get_detected_chips(actualMatch, &nr)
		if chip == nil {
			return chips
		}

		chips = append(chips, &ChipName{*chip, false})
	}
}

func SetParseError(f func(err string, lineno int)) {
	parseError = f
	C.set_sensors_parse_error()
}

func SetParseErrorWfn(f func(err, filename string, lineno int)) {
	parseErrorWfn = f
	C.set_sensors_parse_error_wfn()
}

func SetFatalError(f func(proc, err string)) {
	fatalError = f
	C.set_sensors_fatal_error()
}
