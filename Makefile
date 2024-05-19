build:
	@go build -C cmd/newspaper -o ../../bin/api

run: build
	@./bin/api

milvus:
	docker run -d --name milvus-standalone \
		--security-opt seccomp:unconfined \
		--env-file .env \
		-v $(shell pwd)/volumes/milvus:/var/lib/milvus \
		-v $(shell pwd)/config/embed.yaml:/milvus/configs/embedEtcd.yaml \
		-p 19530:19530 \
		-p 9091:9091 \
		-p 2379:2379 \
		milvus-standalone

test:
	@go test -v ./...
