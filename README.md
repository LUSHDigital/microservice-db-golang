 [![Build Status](https://travis-ci.org/LUSHDigital/microservice-db-golang.svg?branch=master)](https://travis-ci.org/LUSHDigital/microservice-db-golang)
 [![](https://godoc.org/github.com/LUSHDigital/microservice-db-golang?status.svg)](https://godoc.org/github.com/LUSHDigital/microservice-db-golang)


Lush Digital - Micro Service Database (Golang)
---

A set of database packages for Golang microservices. Best used in conjunction
 with [microservice-core-golang](https://github.com/LUSHDigital/microservice-core-golang)

# Package Contents

### Migrator
Migrator utilises the package `github.com/` for processing time-stamped `.sql`
files to be ran as sequenntial migrations against a database. Currently it 
supports the  following databases:
* [CockroachDB](https://www.cockroachlabs.com/product/cockroachdb/)
   
   
# Installation
Install the package as normal, or use your preferred vendoring tool:

`go get -u github.com/LUSHDigital/microservice-db-golang`