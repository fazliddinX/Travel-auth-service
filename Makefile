CURRENT_DIR=$(shell pwd)

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}
exp:
	export DBURL='postgres://postgres:123321@localhost:5432/auth_service?sslmode=disable'

mig-up:
	migrate -path db/migrations -database 'postgres://postgres:123321@localhost:5432/auth_service?sslmode=disable' -verbose up

mig-down:
	migrate -path db/migrations -database 'postgres://postgres:123321@localhost:5432/auth_service?sslmode=disable' -verbose down


mig-create:
	migrate create -ext sql -dir db/migrations -seq auth_service

mig-insert:
	migrate create -ext sql -dir db/migrations -seq follower_table
#   rm -r db/migrations