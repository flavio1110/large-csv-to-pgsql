create-db:
	docker-compose up -d

delete-db:
	docker compose down

gen-file:
	go run . gen

import-stream:
	GOMEMLIMIT=100MiB go run . import-stream

import-read-all:
	GOMEMLIMIT=100MiB go run . import-read-all

import-both:
	GOMEMLIMIT=100MiB go run . import-stream
	GOMEMLIMIT=100MiB go run . import-read-all

benchmark:
	GOMEMLIMIT=100MiB go test -bench=. -benchmem -count 1 -run=^#
