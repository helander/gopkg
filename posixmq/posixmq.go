package posixmq

/*
#cgo LDFLAGS: -lrt

#include <stdlib.h>
#include <signal.h>
#include <fcntl.h>
#include <mqueue.h>

mqd_t mq_open4(const char *name, int oflag, int mode, struct mq_attr *attr) {
	return mq_open(name, oflag, mode, attr);
}
*/
import "C"
import (
//	"fmt"
//	"time"
	"unsafe"
)


type MessageQueueAttribute struct {
	Flags   int
	MaxMsg  int
	MsgSize int
}

const (
	O_RDONLY = C.O_RDONLY
	O_WRONLY = C.O_WRONLY
	O_CREAT  = C.O_CREAT

//	MSGSIZE_MAX     = 16777216
//	MSGSIZE_DEFAULT = MSGSIZE_MAX
)

var (
//	MemoryAllocationError = fmt.Errorf("Memory Allocation Error")
)


func Open(name string, oflag int, mode int, attr *MessageQueueAttribute) (int, error) {
	var cAttr *C.struct_mq_attr
	if attr != nil {
		cAttr = &C.struct_mq_attr{
			mq_flags:   C.long(attr.Flags),
			mq_maxmsg:  C.long(attr.MaxMsg),
			mq_msgsize: C.long(attr.MsgSize),
		}
	}

	h, err := C.mq_open4(C.CString(name), C.int(oflag), C.int(mode), cAttr)
	if err != nil {
		return 0, err
	}

	return int(h), nil
}

func Send(h int, data []byte, priority uint) (int, error) {
	byteStr := *(*string)(unsafe.Pointer(&data))
	rv, err := C.mq_send(C.int(h), C.CString(byteStr), C.size_t(len(data)), C.uint(priority))
	return int(rv), err
}

func Receive(h int, data []byte, len int) (int, error) {
	byteStr := *(*string)(unsafe.Pointer(&data))
	rv, err := C.mq_receive(C.int(h), C.CString(byteStr), C.size_t(len), nil)
	return int(rv), err
}

func Close(h int) (int, error) {
	rv, err := C.mq_close(C.int(h))
	return int(rv), err
}

