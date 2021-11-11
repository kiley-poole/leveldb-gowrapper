package leveldbwrap

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
	"strings"
	"unsafe"
)

type LDB struct {
	DB *C.leveldb_t
}

type Iterator struct {
	Iter *C.leveldb_iterator_t
}

type Options struct {
	opt  *C.leveldb_options_t
	wopt *C.leveldb_writeoptions_t
	ropt *C.leveldb_readoptions_t
}

func (db *LDB) Put(opt *Options, key, val string) {

	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.CString(val)
	defer C.free(unsafe.Pointer(v))

	lenK := C.size_t(len(key))
	lenV := C.size_t(len(val))

	var err *C.char

	C.leveldb_put(db.DB, opt.wopt, k, lenK, v, lenV, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}
}

func (db *LDB) Get(opt *Options, key string) string {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	lenK := C.size_t(len(key))

	var size C.size_t
	var err *C.char

	str1 := C.leveldb_get(db.DB, opt.ropt, k, lenK, &size, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}

	val := C.GoString(str1)

	return val
}

func (db *LDB) Delete(opt *Options, key string) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	lenK := C.size_t(len(key))

	var err *C.char

	C.leveldb_delete(db.DB, opt.wopt, k, lenK, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}
}

func (iter *Iterator) IterateDatabase() {
	C.leveldb_iter_seek_to_first(iter.Iter)
	validIter := C.leveldb_iter_valid(iter.Iter)
	var lenK C.size_t
	var lenV C.size_t
	fmt.Println("The DATABASE CONTENTS!")
	fmt.Println("----------------------")
	for validIter > 0 {
		key := strings.TrimSpace(C.GoString(C.leveldb_iter_key(iter.Iter, &lenK)))
		val := strings.TrimSpace(C.GoString(C.leveldb_iter_value(iter.Iter, &lenV)))
		fmt.Printf("The Key is: %s and the value is: %s\n", key, val)
		C.leveldb_iter_next(iter.Iter)
		validIter = C.leveldb_iter_valid(iter.Iter)
	}
}

func (db *LDB) Close() {
	C.leveldb_close(db.DB)
}

func (opt *Options) DestroyOptions() {
	C.leveldb_options_destroy(opt.opt)
	C.leveldb_writeoptions_destroy(opt.wopt)
	C.leveldb_readoptions_destroy(opt.ropt)
}

func (iter *Iterator) DestroyIterator() {
	C.leveldb_iter_destroy(iter.Iter)
}

func BuildDB(opt *Options, location string) LDB {
	dbname := C.CString(location)
	defer C.free(unsafe.Pointer(dbname))

	var err *C.char

	db := C.leveldb_open(opt.opt, dbname, &err)
	if err != nil {
		panic("err: " + C.GoString(err))
	}

	var leveldb = LDB{DB: db}
	return leveldb

}

func BuildIter(db *LDB, opt *Options) Iterator {

	var iter = Iterator{Iter: C.leveldb_create_iterator(db.DB, opt.ropt)}
	return iter
}

func BuildOptions() Options {
	var opt = Options{
		opt:  C.leveldb_options_create(),
		wopt: C.leveldb_writeoptions_create(),
		ropt: C.leveldb_readoptions_create(),
	}
	C.leveldb_options_set_create_if_missing(opt.opt, 1)
	return opt
}
