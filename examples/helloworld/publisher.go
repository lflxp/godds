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

	"github.com/lflxp/godds"
)

func main() {
	participant, err := godds.DDS_create_participant(godds.DDS_DOMAIN_DEFAULT, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(participant, reflect.TypeOf(participant))
	fmt.Println(godds.DDS_DOMAIN_DEFAULT, reflect.TypeOf(godds.DDS_DOMAIN_DEFAULT))
	var msg C.HelloWorldData_Msg

	topic, err := godds.DDS_create_topic(participant, unsafe.Pointer(&C.HelloWorldData_Msg_desc), "HelloWorldData_Msg", nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("topic", topic, reflect.TypeOf(topic))

	writer, err := godds.DDS_create_writer(participant, topic, nil, nil)
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
