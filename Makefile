
.PHONY: run

PIPELINE_LABEL?=local

run:
	go build -o bin/notes . && ./bin/notes