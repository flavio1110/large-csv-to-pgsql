gen-file:
	go run . gen

import-stream:
	GOMEMLIMIT=100MiB GOGC=100 go run . import-stream

import-read-all:
	GOMEMLIMIT=100MiB GOGC=100 go run . import-read-all

import-both:
	GOMEMLIMIT=100MiB GOGC=100 go run . import-stream
	GOMEMLIMIT=100MiB GOGC=100 go run . import-read-all