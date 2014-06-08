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

/* This program contains the HTTP Server and the Request Handler 
    that calls the appropriate CRUD methods on the DBManager 
*/

package main 

import (
	"net/http" 
	"net/url"
	"encoding/json"
	"fmt"
	"mongodb"
	"strings"
	"strconv"
	"errors"
	"labix.org/v2/mgo/bson"
	"io/ioutil"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	MAX_ROWS = 2000
)

var dbMgr *mongodb.DBManager

func main() {

	 initDB()
	 
	 mux := http.NewServeMux()
	
	 mux.HandleFunc("/gomongo/services/", handleRequest)
     http.ListenAndServe(":8088", mux)
}

func initDB() {

	dbMgr  = mongodb.NewDBManager()
	
	err := dbMgr.InitSession("mongodb://amsAppUser:password@localhost,localhost/AMS")  //TODO: Externalize the URL
	
	if err != nil {
          fmt.Println("DBManager init failed",err)
    } else {
     	  fmt.Println("DBManager Initialized")
    }
}

func handleRequest(response http.ResponseWriter, request *http.Request) {
	
	switch request.Method {

		case GET:
		 	processGET(response,request)  
		case POST:
			processPOST(response,request)
		case PUT:
			processPUT(response,request)
		case DELETE:
			processDELETE(response,request)
	}
}

func processGET(response http.ResponseWriter, request *http.Request) {

	var result []mongodb.Document
	var err error
	var content []byte
	var queryParms mongodb.Document
	var selectParms mongodb.Document
	var docId string
	
	serviceURL := request.URL
	
	collection := getCollectionName(serviceURL)
	
	docId = getDocId(serviceURL)
	
	if len(docId) != 0  { //get document by Id

		result,err = dbMgr.GetDocument(collection,docId)
		
	} else { //run the query using query aprms
	
		queryParms,err = getParms(serviceURL,"q")
		if err == nil { //query parms look good .. 
			selectParms,err = getParms(serviceURL,"select")
		}

		sortParms := getSortParms(serviceURL) //get any sort parms
	
		limit := getLimit(serviceURL)  //get any limit that's passed
	
		if err == nil {
			result,err = dbMgr.RunQuery(collection,queryParms,selectParms,sortParms,limit)
		}
	}

	//check for results
	if err != nil {
		mwError := bson.M{"errorCode":"500","errorMessage":err.Error()}
		content,err = json.MarshalIndent(mwError, "", "  ")
	}else {
		if result == nil || len(result) == 0 {
			content,err = json.MarshalIndent(map[string]mongodb.FieldValue{"errorCode":"404","errorMessage":"No Data Found"}, "", "  ")
		} else {
			content,err = json.MarshalIndent(result, "", "  ")	
		}
	}
	
	//something went wrong marshalling response
	if err != nil { 
		mwError := bson.M{"errorCode":"500","errorMessage":err.Error()}
		content,_ = json.MarshalIndent(mwError, "", "  ")
	}
	
	response.Header().Add("Content-Type","application/json") 

   	response.Write(content)   
}

func processPOST(response http.ResponseWriter, request *http.Request) {

	var err error
	var data mongodb.Document
	var content []byte
	var docId bson.ObjectId
	
	serviceURL := request.URL
	collection := getCollectionName(serviceURL)

	defer request.Body.Close()	
		
	//get POSTed data and unmarshall to JSON	
   	body, _ := ioutil.ReadAll(request.Body)
   	err = json.Unmarshal(body, &data)
    
    //check JSON validity of posteed data 
    if err != nil {
		err = errors.New("Invalid JSON Data: " + err.Error())
	}else {
		err,docId = dbMgr.InsertDocument(collection,data)	
	}
    
	//check for result of the Insert
	if err != nil {
		mwError := bson.M{"errorCode":"500","errorMessage":errors.New("POST Failed - " + err.Error()).Error()}
		content,err = json.MarshalIndent(mwError, "", "  ")
	}else {
		successData := bson.M{"success":true,"message":"Posted Successfully!","_id":docId}
		content,err = json.MarshalIndent(successData, "", "  ")
	}
	
	//something went wrong marshalling response
	if err != nil { 
		mwError := bson.M{"errorCode":"500","errorMessage":err.Error()}
		content,_ = json.MarshalIndent(mwError, "", "  ")
	}
	
	response.Header().Add("Content-Type","application/json") 

   	response.Write(content)   
}

func processPUT(response http.ResponseWriter, request *http.Request) {

	var err error
	var content []byte
	var docId string
	var data mongodb.Document
	
	serviceURL := request.URL
	
	collection := getCollectionName(serviceURL)
	
	docId = getDocId(serviceURL)
	
	if len(docId) != 0  { //delete document
		if bson.IsObjectIdHex(docId) {
			//get POSTed data and unmarshall to JSON	
   			body, _ := ioutil.ReadAll(request.Body)
   			err = json.Unmarshal(body, &data)
    		//check JSON validity of posteed data 
    		if err != nil {
				err = errors.New("Invalid JSON Data: " + err.Error())
			}else {
				query := mongodb.Document{"_id": bson.ObjectIdHex(docId)}
				err = dbMgr.UpdateDocument(collection,query,data)
			}
    	}else {
     		err = errors.New("Invalid Document Id")
    	}
	} else { 
		err = errors.New("Please provide the Id of the document to be Updated")
	}

	//check for result of the Update
	if err != nil {
		mwError := bson.M{"errorCode":"500","errorMessage":errors.New("UPDATE Failed - " + err.Error()).Error()}
		content,err = json.MarshalIndent(mwError, "", "  ")
	}else {
		successData := bson.M{"success":true,"message":"Updated Successfully!"}
		content,err = json.MarshalIndent(successData, "", "  ")
	}
	
	//something went wrong marshalling response
	if err != nil { 
		mwError := bson.M{"errorCode":"500","errorMessage":err.Error()}
		content,_ = json.MarshalIndent(mwError, "", "  ")
	}
	
	response.Header().Add("Content-Type","application/json") 

   	response.Write(content)   
}

func processDELETE(response http.ResponseWriter, request *http.Request) {

	var err error
	var content []byte
	var docId string
	
	serviceURL := request.URL
	
	collection := getCollectionName(serviceURL)
	
	docId = getDocId(serviceURL)
	
	if len(docId) != 0  { //delete document
		if bson.IsObjectIdHex(docId) {
			err = dbMgr.DeleteDocument(collection,mongodb.Document{"_id": bson.ObjectIdHex(docId)})
    	}else {
     		err = errors.New("Invalid Document Id")
    	}
	} else { 
		err = errors.New("Please provide the Id of the document to be Deleted")
	}

	//check for result of the Insert
	if err != nil {
		mwError := bson.M{"errorCode":"500","errorMessage":errors.New("DELETE Failed - " + err.Error()).Error()}
		content,err = json.MarshalIndent(mwError, "", "  ")
	}else {
		successData := bson.M{"success":true,"message":"Deleted Successfully!"}
		content,err = json.MarshalIndent(successData, "", "  ")
	}
	
	//something went wrong marshalling response
	if err != nil { 
		mwError := bson.M{"errorCode":"500","errorMessage":err.Error()}
		content,_ = json.MarshalIndent(mwError, "", "  ")
	}
	
	response.Header().Add("Content-Type","application/json") 

   	response.Write(content)   
}

func getCollectionName(serviceURL *url.URL) (string) {

	path := serviceURL.Path    // 'gomongo/services/collectionName/documentId
	
	pathSplits := strings.Split(path,"/")

	return pathSplits[3]
}

func getDocId(serviceURL *url.URL) (string) {

	path := serviceURL.Path    // 'gomongo/services/collectionName/documentId
	
	pathSplits := strings.Split(path,"/")
	
	if len(pathSplits) > 4 {
		return pathSplits[4]
	}
	
	return ""
}

func getParms(serviceURL *url.URL,parmType string) (parmsDocument mongodb.Document,err error ) {

	var parmsString string 
	
	rawQuery := serviceURL.RawQuery
	
    parmsMap, _ := url.ParseQuery(rawQuery)

    if(parmsMap[parmType] != nil && len(parmsMap[parmType][0]) > 0) {
   		parmsString =  parmsMap[parmType][0]
   		err = json.Unmarshal([]byte(parmsString) , &parmsDocument)
   	}else {
   		return nil,nil
   	}
   	
   	if err != nil {
   		if parmType == "q" {
   			err = errors.New("Invalid Query Parms- " + err.Error())
   		} else {
   			err = errors.New("Invalid Select Parms-" + err.Error())
   		}
   	}
 
  	return parmsDocument,err
}

func getSortParms(serviceURL *url.URL) (sortParms []string) {
 
	rawQuery := serviceURL.RawQuery
	
    parmsMap, _ := url.ParseQuery(rawQuery)

    if(parmsMap["sort"] != nil && len(parmsMap["sort"][0]) > 0) {
   		sortParms =  parmsMap["sort"]
   	}
 
  	return sortParms
}

func getLimit(serviceURL *url.URL) (limit int) {

	var err error
	limit = MAX_ROWS
	
	rawQuery := serviceURL.RawQuery
	
    parmsMap, _ := url.ParseQuery(rawQuery)

    if(parmsMap["limit"] != nil && len(parmsMap["limit"][0]) > 0) {
    	limitString := parmsMap["limit"][0]
   		limit,err = strconv.Atoi(limitString) 
   		if err != nil {
   			limit = MAX_ROWS
   		}
   	}
 
  	return limit
}

