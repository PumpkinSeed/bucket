# Bucket

[![Build Status](https://travis-ci.com/PumpkinSeed/bucket.svg?branch=master)](https://travis-ci.com/PumpkinSeed/bucket)
[![Go Report Card](https://goreportcard.com/badge/github.com/PumpkinSeed/bucket)](https://goreportcard.com/report/github.com/PumpkinSeed/bucket)
[![Documentation](https://godoc.org/github.com/PumpkinSeed/bucket?status.svg)](http://godoc.org/github.com/PumpkinSeed/bucket)

Simply Couchbase framework.

Project specifically focus on the one bucket as database approach and makes it easier to manage complex data sets. It tries to get rid of the embedded jsons per document and separate them into different documents behind the scenes.

##### Features:
- Automatic index generator with indexable tags.
- Simple usage through the handler.
- Following best practices of Couchbase usage behind the scenes, which doesn't affect the user of the library.

##### How to use:

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

You can also index structs adding the `indexable="true"` tag to the struct, and then calling `bucket.Index({context.Context}, yourStruct)`.