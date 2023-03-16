package main

/*
#cgo LDFLAGS: -L ../../library/lib -lddsc ${SRCDIR}/HelloWorldData.o
#cgo CFLAGS: -I ../../library/include
#include "HelloWorldData.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "dds/dds.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"time"
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

// func ErrorCheck(err C.dds_entity_t, flags uint8, where string) {
// 	C.dds_err_check(err, C.uint(flags), C.CString(where))
// }

type DDS_entity_t C.dds_entity_t
type DDS_return_t C.dds_return_t
type DomainID C.dds_domainid_t
type Listener C.dds_listener_t

type QoS C.dds_qos_t

func DDS_create_participant(domainID DomainID, qos *QoS, listener *Listener) (DDS_entity_t, error) {
	participant := C.dds_create_participant((C.dds_domainid_t)(domainID), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	// C.printf(participant)
	// fmt.Println(fmt.Sprintf("%v", participant), reflect.TypeOf(participant))
	// if participant < 0 {
	// 	fmt.Println("error")
	// } else {
	// 	fmt.Println("ok")
	// }
	if participant < 0 {
		return (DDS_entity_t)(participant), CddsErrorType(participant)
	}
	return (DDS_entity_t)(participant), nil
}

func DDS_create_topic(participant DDS_entity_t, desc unsafe.Pointer, name string, qos *QoS, listener *Listener) (DDS_entity_t, error) {
	topic := C.dds_create_topic((C.dds_entity_t)(participant), (*C.dds_topic_descriptor_t)(desc), C.CString(name), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if topic < 0 {
		return (DDS_entity_t)(topic), CddsErrorType(topic)
	}
	return (DDS_entity_t)(topic), nil
}

func DDS_create_writer(participant DDS_entity_t, topic DDS_entity_t, qos *QoS, listener *Listener) (DDS_entity_t, error) {
	writer := C.dds_create_writer((C.dds_entity_t)(participant), (C.dds_entity_t)(topic), (*C.dds_qos_t)(qos), (*C.dds_listener_t)(listener))
	if writer < 0 {
		return (DDS_entity_t)(writer), CddsErrorType(writer)
	}
	return (DDS_entity_t)(writer), nil
}

// func DDS_
func DDS_write(writer DDS_entity_t, msg interface{}) error {
	rc := C.dds_write((C.int)(writer), unsafe.Pointer(&msg))
	fmt.Println(rc)
	return nil
}

func main() {
	participant, err := DDS_create_participant(DDS_DOMAIN_DEFAULT, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(participant, reflect.TypeOf(participant))
	fmt.Println(DDS_DOMAIN_DEFAULT, reflect.TypeOf(DDS_DOMAIN_DEFAULT))
	var msg C.HelloWorldData_Msg

	topic, err := DDS_create_topic(participant, unsafe.Pointer(&C.HelloWorldData_Msg_desc), "HelloWorldData_Msg", nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("topic", topic, reflect.TypeOf(topic))

	writer, err := DDS_create_writer(participant, topic, nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("writer", writer, reflect.TypeOf(writer))
	fmt.Println("=== [Publisher]  Waiting for a reader to be discovered ...\n")
	C.fflush(C.stdout)

	rc := C.dds_set_status_mask((C.int)(writer), C.DDS_PUBLICATION_MATCHED_STATUS)
	fmt.Println(rc)
	if rc != C.DDS_RETCODE_OK {
		panic(fmt.Sprintf("dds_set_status_mask:%s\n", C.dds_strretcode(-rc)))
	}

	var status = (C.uint32_t)(0)
	for status > 0 && C.DDS_PUBLICATION_MATCHED_STATUS > 0 {
		rc = C.dds_get_status_changes((C.int)(writer), &status)
		fmt.Println("status: rc", rc)
		if rc != C.DDS_RETCODE_OK {
			panic(fmt.Sprintf("dds_get_status_changes:%s\n", C.dds_strretcode(-rc)))
		}
		// C.dds_sleepfor(C.DDS_MSECS(20))
		time.Sleep(time.Microsecond * 20)
	}

	for x := 0; x < 100000; x++ {
		ms := C.CString("Hello, world!")
		msg.userID = (C.int)(x)
		msg.message = (*C.char)(unsafe.Pointer(ms))

		fmt.Println("=== [Publisher] Writing :")
		fmt.Printf("Message %v %v\n", msg.userID, C.GoString(msg.message))
		C.fflush(C.stdout)

		// err = DDS_write(writer, &msg)
		// if err != nil {
		// 	panic(err)
		// }

		rc = C.dds_write((C.int)(writer), unsafe.Pointer(&msg))
		fmt.Println(rc, C.DDS_RETCODE_OK)
		if rc != C.DDS_RETCODE_OK {
			panic(fmt.Sprintf("dds_write: %s\n", C.dds_strretcode(-rc)))
		}
		time.Sleep(time.Second)
	}

	rc = C.dds_delete((C.int)(participant))
	if rc != C.DDS_RETCODE_OK {
		panic(fmt.Sprintf("dds_delete: %s\n", C.dds_strretcode(-rc)))
	}
	fmt.Println(C.EXIT_SUCCESS)
}
