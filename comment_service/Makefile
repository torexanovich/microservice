run:
	go run cmd/main.go

proto-gen:
	./scripts/gen-proto.sh

migrate_up:
	migrate -path migrations -database postgres://smn:123@localhost:5434/smn up

migrate_down:
	migrate -path migrations -database postgres://smn:123@localhost:5434/smn down

migrate_force:
	migrate -path migrations -database postgres://smn:123@localhost:5434/smn -verbose force 0

.PHONY: start migrateup migratedown