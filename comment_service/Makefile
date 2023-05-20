run:
	go run cmd/main.go

proto-gen:
	./scripts/gen-proto.sh

migrate_up:
	migrate -path migrations -database postgres://smn:123@localhost:5434/human up

migrate_down:
	migrate -path migrations -database postgres://smn:123@localhost:5434/human down

migrate_force:
	migrate -path migrations -database postgres://smn:123@localhost:5434/human -verbose force 0

.PHONY: start migrateup migratedown