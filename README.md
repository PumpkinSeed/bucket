# Bucket

[![Build Status](https://travis-ci.com/PumpkinSeed/bucket.svg?branch=master)](https://travis-ci.com/PumpkinSeed/bucket)
[![Go Report Card](https://goreportcard.com/badge/github.com/PumpkinSeed/bucket)](https://goreportcard.com/report/github.com/PumpkinSeed/bucket)
[![Documentation](https://godoc.org/github.com/PumpkinSeed/bucket?status.svg)](http://godoc.org/github.com/PumpkinSeed/bucket)

Simple Couchbase framework.

Project specifically focuses on the one bucket as database approach, and makes it easier to manage complex data sets. It tries to get rid of the embedded jsons per document and separates them into different documents behind the scene.

### Disclaimer:
**DO NOT USE IN PRODUCTION.** This is still a work in progress. We will not take responsibility for any breaks in the code that happen after a new version comes out.

#### Features:
- Automatic index generator with indexable tags.
- Simple usage through the handler.
- Following best practices of Couchbase usage behind the scene, which doesn't affect the user of the library.

#### Rules:

1. Only struct can be referenced

#### How to use:

Create a new handler with the New function, that needs a configuration.
```go
bucket.New( {bucket.Configuration} )

type Configuration struct {
    // The address of the couchbase server
	ConnectionString string 

    // Username and password to access couchbase
	Username         string 
	Password         string 
	
    // The name and password of the bucket you want to use
	BucketName       string 
	BucketPassword   string 

    // The separator of your choice, this will separate the prefix from the field name
	Separator        string 
}
```

After that you can use the Insert, Get, Remove, Upsert, Touch, GetAndTouch and Ping methods of the handler.

```go
package main

import (
    "context"
    "fmt"

    "github.com/PumpkinSeed/bucket"
    "github.com/couchbase/gocb"
)

var conf = &bucket.Configuration{
    ConnectionString: "myServer.org:1234",
    Username: "cbUsername",
    Password: "cbPassword",
    BucketName: "testBucket",
    BucketPassword: "testBucketPassword",
    Separator: "::",
}

type myStruct struct {
    justAField string `json:"just_a_field"`
}

func main() {
    var in = &myStruct{
        justAField: "basic",
    }
    var cas bucket.Cas
    var typ = "prefix"
    ctx := context.Background()

    h, err := bucket.New(conf)
    if err != nil {
        // handle error
    }
    
    // insert
    cas, id, err := h.Insert(ctx, typ, "myID", in, 0)
    if err != nil {
        // handle error
    }

    // get
    var out = &myStruct{}
    err = h.Get(ctx, typ, id, out)
    if err != nil {
        // handle error
    }

    // touch
    err = h.Touch(ctx, typ, id, in, 0)
    if err != nil {
        // handle error
    }

    // get and touch
    var secondOut = &myStruct{}
    err = h.GetAndTouch(ctx, typ, id, secondOut, 0)
    if err != nil {
        // handle error
    }

    // ping
    var services []gocb.ServiceType
    res, err := h.Ping(ctx, services)
    if err != nil {
        // handle error
    }

    fmt.Println(res)

    // upsert
    in.justAField = "updated"
    cas, newID, err := h.Upsert(ctx, typ, id, in, 0)

    // remove
    err = h.Remove(ctx, typ, newID, in)
    if err != nil {
        // handle error
    }
}
```

**Important:** 
- The typ parameter will be the prefix of the initial struct, so you should use the same value for the same types!
- IDs should be unique, if the parameter is an empty string (`""`) a globally unique ID will be automatically generated!

#### Additional:

Embedded structs can be separated into a a new entry with the `cb_referenced` tag. The value will decide the typ of the struct.
```go
type example struct {
    refMe       *refMe      `json:"ref_me" cb_referenced:"separate_entry"`
    ignoreMe    *ignoreMe   `json:"ignore_me"`
}

type refMe struct {
    referencedField int `json:"referenced_field"`
}

type ignoreMe struct {
    notReferencedField int `json:"not_referenced_field"`
}
```

You can also index structs adding the `cb_indexable:"true"` tag to the field, and then calling `bucket.Index({context.Context}, yourStruct)`.
```go
type example struct {
    indexMe         string `json:"index_me" cb_indexable:"true"`
    butNotThisOne   string `json:"but_not_this_one"`
}
```