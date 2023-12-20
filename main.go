package main

import (
	"fmt"
)

var dbName = "trans_test"
var collName = "coll1"

// NOTE:  (IllegalOperation) Transaction numbers are only allowed on a replica set member or mongos

var connString = "mongodb://localhost:27017"

func main() {
	fmt.Println(dbName, collName)

	collections := []string{"coll1", "coll2", "coll3"}
	// test transaction
	TestMultiInsertTransactionCommit(dbName, collections...)
}
