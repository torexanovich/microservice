run:
	go run cmd/main.go

proto-gen:
	./scripts/gen-proto.sh

migrate_up:
	migrate -path migrations -database postgres://smn:123@localhost:5432/postdb up

migrate_down:
	migrate -path migrations -database postgres://smn:123@localhost:5432/postdb down

migrate_force:
	migrate -path migrations -database postgres://smn:123@localhost:5432/postdb -verbose force 0

.PHONY: start migrateup migratedown