package main

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lleveldb
#include "leveldb/c.h"
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type ldb struct {
	db *C.leveldb_t
}

type opt struct {
	opt  *C.leveldb_options_t
	wopt *C.leveldb_writeoptions_t
	ropt *C.leveldb_readoptions_t
}

func (db *ldb) Put(opt *opt, key, val string) {

	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.CString(val)
	defer C.free(unsafe.Pointer(v))

	lenK := C.size_t(len(key))
	lenV := C.size_t(len(val))

	var err *C.char

	C.leveldb_put(db.db, opt.wopt, k, lenK, v, lenV, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}
}

func (db *ldb) Get(opt *opt, key string) string {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	lenK := C.size_t(len(key))

	var size C.size_t
	var err *C.char

	str1 := C.leveldb_get(db.db, opt.ropt, k, lenK, &size, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}

	val := C.GoString(str1)

	return val
}

func open(opt *opt, location string) *C.leveldb_t {
	dbname := C.CString(location)
	defer C.free(unsafe.Pointer(dbname))

	var err *C.char

	db := C.leveldb_open(opt.opt, dbname, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}

	return db
}

func buildOptions() *opt {
	var opt = opt{
		opt:  C.leveldb_options_create(),
		wopt: C.leveldb_writeoptions_create(),
		ropt: C.leveldb_readoptions_create()}
	C.leveldb_options_set_create_if_missing(opt.opt, 1)
	return &opt
}

func main() {
	opt := buildOptions()
	levelDB := ldb{db: open(opt, "/tmp/newDB")}
	levelDB.Put(opt, "Test1", "Value1")
	fmt.Println(levelDB.Get(opt, "Test1"))
}
