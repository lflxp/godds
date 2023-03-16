package godds

// #cgo LDFLAGS: -L {SRCDIR}/library/lib
// #cgo CFLAGS: -I {SRCDIR}/library/include

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "dds/dds.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type CddsErrorType uint16

const (
	DDS_DOMAIN_DEFAULT DomainID      = (DomainID)(0)
	Ok                 CddsErrorType = iota
	Error
	Unsupported
	BadParameter
	PreconditionNotMet
	OutOfResource
	NotEnabled
	ImmutablePolicy
	InconsistencyPolicy
	AlreadyDeleted
	TimeOut
	NoData
	IllegalOperation
)

func (e CddsErrorType) Error() string {
	return []string{
		"Success",
		"Non specific error",
		"Feature unsupported",
		"Bad parameter value",
		"Precondition for operation not met",
		"Out of resources",
		"Configurable feature is not enabled",
		"Attempt is made to modify an immutable policy",
		"Policy is used with inconsistent values",
		"Attempt is made to delete something more than once",
		"Timeout",
		"Expected data is not provided",
		"Function is called when it should not be",
		"credentials are not enough to use the function",
	}[int(e)]
}

// func ErrorCheck(err C.DDS_ENTITY_T, flags uint8, where string) {
// 	C.dds_err_check(err, C.uint(flags), C.CString(where))
// }

type DDS_ENTITY_T C.dds_entity_t
type DDS_RETURN_T C.dds_return_t
type DomainID C.dds_domainid_t
type Listener C.dds_listener_t
type Infos C.dds_sample_info_t

type QoS C.dds_qos_t

func DDS_create_participant(domainID DomainID, qos *QoS, listener *Listener) (DDS_ENTITY_T, error) {
	participant := C.dds_create_participant((C.dds_domainid_t)(domainID), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	// C.printf(participant)
	// fmt.Println(fmt.Sprintf("%v", participant), reflect.TypeOf(participant))
	// if participant < 0 {
	// 	fmt.Println("error")
	// } else {
	// 	fmt.Println("ok")
	// }
	if participant < 0 {
		return (DDS_ENTITY_T)(participant), CddsErrorType(participant)
	}
	return (DDS_ENTITY_T)(participant), nil
}

func DDS_create_topic(participant DDS_ENTITY_T, desc unsafe.Pointer, name string, qos *QoS, listener *Listener) (DDS_ENTITY_T, error) {
	topic := C.dds_create_topic((C.dds_entity_t)(participant), (*C.dds_topic_descriptor_t)(desc), C.CString(name), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if topic < 0 {
		return (DDS_ENTITY_T)(topic), CddsErrorType(topic)
	}
	return (DDS_ENTITY_T)(topic), nil
}

func DDS_create_writer(participant DDS_ENTITY_T, topic DDS_ENTITY_T, qos *QoS, listener *Listener) (DDS_ENTITY_T, error) {
	writer := C.dds_create_writer((C.dds_entity_t)(participant), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if writer < 0 {
		return (DDS_ENTITY_T)(writer), CddsErrorType(writer)
	}
	return (DDS_ENTITY_T)(writer), nil
}

func DDS_create_reader(participant DDS_ENTITY_T, topic DDS_ENTITY_T, qos *QoS, listener *Listener) (DDS_ENTITY_T, error) {
	reader := C.dds_create_reader((C.dds_entity_t)(participant), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if reader < 0 {
		return (DDS_ENTITY_T)(reader), CddsErrorType(reader)
	}
	return (DDS_ENTITY_T)(reader), nil
}

// func DDS_
func DDS_write(writer DDS_ENTITY_T, msg interface{}) error {
	rc := C.dds_write((C.int)(writer), unsafe.Pointer(&msg))
	fmt.Println(rc)
	return nil
}
