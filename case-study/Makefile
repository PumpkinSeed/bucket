protogen:
	protoc --go_out=plugins=grpc:. models/company.proto
	protoc-go-inject-tag -input=./models/company.pb.go