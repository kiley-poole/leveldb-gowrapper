package main

import (
	"fmt"

	"github.com/kiley-poole/leveldb-gowrapper/leveldbwrap"
)

func getDBName() string {
	var db string
	fmt.Println("Enter the database name you'd like to access: ")
	fmt.Scanln(&db)
	return db
}

func displayAndGetOption() int {
	var opt int
	fmt.Println("\nEnter 1 to add a new value to the database")
	fmt.Println("Enter 2 to retrieve a value from the database")
	fmt.Println("Enter 3 to delete a value from the database")
	fmt.Println("Enter 4 to print all key value pairs in the database")
	fmt.Println("Enter 0 to exit")
	fmt.Println("Select Your Option: ")
	fmt.Scanln(&opt)
	return opt
}

func addVal(db *leveldbwrap.LDB, opt *leveldbwrap.Options) {
	var key, val string
	fmt.Println("\nEnter the new key: ")
	fmt.Scanln(&key)
	fmt.Println("Enter the new val: ")
	fmt.Scanln(&val)
	db.Put(opt, key, val)
}

func getVal(db *leveldbwrap.LDB, opt *leveldbwrap.Options) {
	var key string
	fmt.Println("\nEnter the key you'd like to retrieve: ")
	fmt.Scanln(&key)
	val := db.Get(opt, key)
	fmt.Printf("\nThe value is: %s\n\n", val)
}

func deleteVal(db *leveldbwrap.LDB, opt *leveldbwrap.Options) {
	var key string
	fmt.Println("\nEnter the key you'd like to delete: ")
	fmt.Scanln(&key)
	db.Delete(opt, key)
	fmt.Printf("\nThe value was deleted")
}

func displayDB(iter *leveldbwrap.Iterator) {
	iter.IterateDatabase()
}

func closeAll(db *leveldbwrap.LDB, opt *leveldbwrap.Options, iter *leveldbwrap.Iterator) {
	iter.DestroyIterator()
	db.Close()
	opt.DestroyOptions()
}

func main() {
	opt := leveldbwrap.BuildOptions()
	dbName := getDBName()
	levelDB := leveldbwrap.BuildDB(&opt, dbName)
	iterator := leveldbwrap.BuildIter(&levelDB, &opt)

Loop:
	for {
		menuOpt := displayAndGetOption()
		switch menuOpt {
		case 1:
			addVal(&levelDB, &opt)
		case 2:
			getVal(&levelDB, &opt)
		case 3:
			deleteVal(&levelDB, &opt)
		case 4:
			displayDB(&iterator)
		case 0:
			closeAll(&levelDB, &opt, &iterator)
			break Loop
		}
	}

}
