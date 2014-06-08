gomongo
=======

A Simple REST API for mongoDB in GO!

Enables the mongoDB CRUD operations on any collection using simple HTTP REST calls. 

This is my first serious attempt at GO programming language and mongoDB. The code works but I am sure there is a lot of room for re-factoring and clean up. 

To use it in a real web app environment, the code needs to be beefed up with some authentication code around the services.

This project consists of the following three files. It has a dependency on labix.org/v2/mgo

- mongoweb.go :  Main program (HTTP Server)- handles GET,POST,PUT and DELETE
- db.go       :  Provides wrapper methods to mGO(a native mongoDB driver in GO) CRUD fucntions
- testdb.go   :  Test functions for db.go

USAGE:

- GET http://localhost:8088/gomongo/services/people/6466761235764 
- GET http://localhost:8088/gomongo/services/people/
- POST http://localhost:8088/gomongo/services/people/ <br>
	Body: {"name":"YTRETRERETYRE","phone":"+55 53 8116 9639"} <br>
- PUT http://localhost:8088/gomongo/services/people/6466761235764 <br>
	Body: {"name":"NEW NAME","phone":"+55 53 8116 9639"}<br>
- DELETE http://localhost:8088/gomongo/services/people/6466761235764 <br>

Note: <br>
	/gomongo/services is the URL that the service listens to 'people' is the name of the mongodb collections in the above examples
	