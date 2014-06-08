/*
	This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    
    Created By: Hari Kolasani. June 8, 2014.

*/

package main 

import (
    "fmt"
    "mongodb"
    "encoding/json"
    "labix.org/v2/mgo/bson"
)

var dbMgr *mongodb.DBManager

func main() {

	initDB() 
	
	docId := testInsert()
    
    queryDocument := mongodb.Document{"_id":docId}
    selectDocument := mongodb.Document{"name":1}
    sortParms := []string{"name"}
    
    testQuery(queryDocument,selectDocument,sortParms,10)
    
    testUpdate(docId)
    	
    testDelete(docId)
    
    //testDelete(bson.ObjectIdHex("538e527aa5f3170fe9000001"))  This Works too!
    
    dbMgr.Term()
}

func initDB() {

	dbMgr  = mongodb.NewDBManager()
	
	err := dbMgr.InitSession("mongodb://localhost,localhost/test")
	
	if err != nil {
          fmt.Println("DBManager init failed",err)
    } else {
     	  fmt.Println("DBManager Initialized")
    }
}

func testQuery(queryDoc mongodb.Document,selectDoc mongodb.Document,sortParms []string,limit int) {

    result,err := dbMgr.RunQuery("people",queryDoc,selectDoc,sortParms,10)
    
    if err != nil {
          fmt.Println("DBManager RunQuery failed",err)
    } else {
     	  content,err1 := json.MarshalIndent(result, "", "  ")
     	  if err1 != nil {
       		   fmt.Println("DBManager Result Marshalling failed",err1)
    	  } else {
          		fmt.Println("Result:", string(content))
          }
    }   
}

func testInsert() (docId bson.ObjectId) {

	person1 := mongodb.Document{"name":"XSDSXSXSDXSXYYY","phone":"+55 53 8116 9639"}
	
    err,docId := dbMgr.InsertDocument("people",person1)
    
    if err != nil {
          fmt.Println("DBManager Insert failed",err)
    } else {
     	  fmt.Println("DBManager Inserted Successfully:",docId)
    }
    
    return docId
}

func testUpdate(docId bson.ObjectId)  {

	query := mongodb.Document{"_id": docId}

	properties := mongodb.Document{"name":"HGFHGFHGFHFHF","phone":"+55 53 8116 9639"}
	
    err := dbMgr.UpdateDocument("people",query,properties)
    
    if err != nil {
          fmt.Println("DBManager Update failed",err)
    } else {
     	  fmt.Println("DBManager Updated Successfully:",docId)
    }
}

func testDelete(docId bson.ObjectId) {

	err := dbMgr.DeleteDocument("people",mongodb.Document{"_id": docId})
    
    if err != nil {
          fmt.Println("DBManager Delete failed",err)
    } else {
     	  fmt.Println("DBManager Deleted Successfully",docId)
    }
}


