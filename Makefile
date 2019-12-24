
.PHONY: run

PIPELINE_LABEL?=local

run:
	go build . && ./notes