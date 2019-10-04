testing:
	PKG_TEST=testing go test ./...

dev:
	docker-compose -f test/docker-compose.yml down
	docker-compose -f test/docker-compose.yml up -d