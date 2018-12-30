package go_linux_sensors

/*
#include "bridge.h"
#cgo LDFLAGS: -lsensors
*/
import "C"

import (
	"reflect"
	"unsafe"
)

type Error struct {
	errnr C.int
}

func (e Error) Error() string {
	{
		errorsLock.RLock()
		err, hasErr := errors[e.errnr]
		errorsLock.RUnlock()

		if hasErr {
			return err
		}
	}

	err := cStr2str(C.sensors_strerror(e.errnr))

	errorsLock.Unlock()
	errors[e.errnr] = err
	errorsLock.Lock()

	return err
}

type ChipName struct {
	chipName    C.struct_sensors_chip_name
	freeOnClose bool
}

func (c *ChipName) MarshalBinary() ([]byte, error) {
	res := make([]byte, int(C.strlen(c.chipName.prefix))+int(C.strlen(c.chipName.path))+256)

	errSSCN := C.sensors_snprintf_chip_name(
		(*C.char)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&res)).Data)),
		C.Tsize(len(res)),
		&c.chipName,
	)
	if errSSCN < C.int(0) {
		return nil, Error{errSSCN}
	}

	return res[:int(errSSCN)], nil
}

func (c *ChipName) GetFeatures() []Feature {
	nr := C.int(0)
	features := []Feature{}

	for {
		feature := C.sensors_get_features(&c.chipName, &nr)
		if feature == nil {
			return features
		}

		features = append(features, Feature{feature})
	}
}

func (c *ChipName) GetAllSubfeatures(feature Feature) []Subfeature {
	nr := C.int(0)
	subfeatures := []Subfeature{}

	for {
		subfeature := C.sensors_get_all_subfeatures(&c.chipName, feature.feature, &nr)
		if subfeature == nil {
			return subfeatures
		}

		subfeatures = append(subfeatures, Subfeature{subfeature})
	}
}

func (c *ChipName) GetSubfeature(feature Feature, typ SubfeatureType) (Subfeature, bool) {
	subfeature := C.sensors_get_subfeature(&c.chipName, feature.feature, C.sensors_subfeature_type(typ))
	if subfeature == nil {
		return Subfeature{}, false
	}

	return Subfeature{subfeature}, true
}

func (c *ChipName) GetLabel(feature Feature) (string, bool) {
	label := C.sensors_get_label(&c.chipName, feature.feature)
	if label == nil {
		return "", false
	}

	defer C.free(unsafe.Pointer(label))

	return cStr2str(label), true
}

func (c *ChipName) GetValue(subfeatNr int) (float64, error) {
	var value C.double

	if errSGV := C.sensors_get_value(&c.chipName, C.int(subfeatNr), &value); errSGV != C.int(0) {
		return 0, Error{errSGV}
	}

	return float64(value), nil
}

func (c *ChipName) SetValue(subfeatNr int, value float64) error {
	if errSSV := C.sensors_set_value(&c.chipName, C.int(subfeatNr), C.double(value)); errSSV != C.int(0) {
		return Error{errSSV}
	}

	return nil
}

func (c *ChipName) DoChipSets() error {
	if errSDCS := C.sensors_do_chip_sets(&c.chipName); errSDCS != C.int(0) {
		return Error{errSDCS}
	}

	return nil
}

func (c *ChipName) Close() error {
	if c.freeOnClose {
		C.sensors_free_chip_name(&c.chipName)
	}

	return nil
}

func (c *ChipName) GetPrefix() string {
	return cStr2str(c.chipName.prefix)
}

func (c *ChipName) GetBus() BusId {
	return BusId{c.chipName.bus}
}

func (c *ChipName) GetAddr() int {
	return int(c.chipName.addr)
}

func (c *ChipName) GetPath() string {
	return cStr2str(c.chipName.path)
}

type BusId struct {
	busId C.struct_sensors_bus_id
}

func (b BusId) GetType() int16 {
	return int16(b.busId._type)
}

func (b BusId) GetNr() int16 {
	return int16(b.busId.nr)
}

func (b BusId) GetAdapterName() (string, bool) {
	adapterName := C.sensors_get_adapter_name(&b.busId)
	if adapterName == nil {
		return "", false
	}

	return cStr2str(adapterName), true
}

type FeatureType int

type Feature struct {
	feature *C.TCsensors_feature
}

func (f Feature) GetName() string {
	return cStr2str(f.feature.name)
}

func (f Feature) GetNumber() int {
	return int(f.feature.number)
}

func (f Feature) GetType() FeatureType {
	return FeatureType(f.feature._type)
}

type SubfeatureType int

type Subfeature struct {
	subfeature *C.TCsensors_subfeature
}

func (s Subfeature) GetName() string {
	return cStr2str(s.subfeature.name)
}

func (s Subfeature) GetNumber() int {
	return int(s.subfeature.number)
}

func (s Subfeature) GetType() SubfeatureType {
	return SubfeatureType(s.subfeature._type)
}

func (s Subfeature) GetMapping() int {
	return int(s.subfeature.mapping)
}

func (s Subfeature) GetFlags() uint {
	return uint(s.subfeature.flags)
}
