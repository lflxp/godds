package main

/*
#cgo LDFLAGS: -L ../../library/lib -lddsc ${SRCDIR}/HelloWorldData.o
#cgo CFLAGS: -I ../../library/include
#include "HelloWorldData.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "dds/dds.h"


#define MAX_SAMPLES 1
void *samples[MAX_SAMPLES];
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

	/* Create a reliable Reader. */
	qos := C.dds_create_qos()
	C.dds_qset_reliability(qos, C.DDS_RELIABILITY_RELIABLE, C.DDS_SECS(10))
	reader, err := godds.DDS_create_reader(participant, topic, qos, nil)
	if err != nil {
		panic(err)
	}
	C.dds_delete_qos(unsafe.Pointer(qos))

	fmt.Println("=== [Subscriber]  Waiting for a sample ...\n")
	C.fflush(C.stdout)

	C.samples[0] = C.HelloWorldData_Msg__alloc()

	for {
		rc := C.dds_read((C.int)(reader), C.samples, godds.Infos[C.MAX_SAMPLES], C.MAX_SAMPLES, C.MAX_SAMPLES)
		fmt.Println("status: rc", rc)
		if rc < 0 {
			panic(fmt.Sprintf("dds_read:%s\n", C.dds_strretcode(-rc)))
		}

		if rc > 0 && godds.Infos[0].vaild_data {
			msg = (*C.HelloWorldData_Msg)(godds.Infos[0])
			fmt.Println("=== [Subscriber] Receoved")
			fmt.Printf("Message (%d, %s)\n", msg.userID, msg.message)
			C.fflush(C.stdout)
			break
		} else {
			time.Sleep(time.Microsecond * 20)
		}

	}

	C.HelloWorldData_Msg_free(godds.samples[0], C.DDS_FREE_ALL)

	rc = C.dds_delete((C.int)(participant))
	if rc != C.DDS_RETCODE_OK {
		panic(fmt.Sprintf("dds_delete: %s\n", C.dds_strretcode(-rc)))
	}
	fmt.Println(C.EXIT_SUCCESS)
}
