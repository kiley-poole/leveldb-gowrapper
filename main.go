package main

/*
#include "leveldb/include/leveldb/c.h"
*/
import "C"

func main() {
	type levelDB struct {
		db  *C.leveldb_t
		opt *C.leveldb_options_t
	}

}
