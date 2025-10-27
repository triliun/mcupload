APP_NAME=Mcupload
APP_DIR=.
TEST_E2E_DIR=./test/e2e/

# PORT?=4000
# DSN?=postgres://postgres:root@localhost:5432/snippetbox?sslmode=disable

.PHONY: run build clean

run:
	@echo "Running aplication..."
	go run $(APP_DIR)
# 	go run $(APP_DIR) -addr=$(PORT) -dsn=$(DSN)

te:
	@echo "Running test E2E..."
	go test -v $(TEST_E2E_DIR)

build:
	@echo "Build application..."
	go buil -o bin/$(APP_NAME) $(APP_DIR)

clean:
	@echo "Cleaning binary file..."
	rm -f bin/$(APP_NAME)

pstart:
	@echo "Start database postgresSQL..."
	systemctl start postgresql.service

pstop:
	@echo "Stop database postgresSQL..."
	systemctl stop postgresql.service

help:
	@echo "make run - Running application (default)"
	@echo "make te - Running test E2E"
	@echo "make build - Build application"
	@echo "make clean - Cleaning binary file"
	@echo "make pstart - Starting databse postgresSQL"
	@echo "make pstop - Stoping databse postgresSQL"
