run:
	go run cmd/main.go

proto-gen:
	./scripts/gen-proto.sh

migrate_up:
	migrate -path migrations -database postgres://smn:123@localhost:5432/userd -verbose up

migrate_down:
	migrate -path migrations -database postgres://smn:123@localhost:5432/userd -verbose down

migrate_force:
	migrate -path migrations -database postgres://smn:123@localhost:5432/userd -verbose force 0

test:
	go test ./...
.PHONY: start migrateup migratedown