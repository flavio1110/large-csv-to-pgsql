gen-file:
	go run . gen

import-stream:
	export GOMEMLIMIT=5MiB
	go build -o ./bin/importer
	./bin/importer import-stream

import-read-all:
	go build -o ./bin/importer
	./bin/importer import-read-all

import-both:
	export GOMEMLIMIT=5MiB
	export GOGC=off
	go build -o ./bin/importer
	./bin/importer import-stream
	./bin/importer import-read-all