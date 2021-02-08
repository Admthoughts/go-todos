
GO=$$(which go)
GO_TEST_FLAGS=-v -cover
DOCKER_DB=todo-list-test-db

test:
	$(GO) test $(GO_TEST_FLAGS)

int_db:
	@echo "running db in docker"
	@docker run -d --rm --name $(DOCKER_DB) \
	-p 5432:5432 \
	-e POSTGRES_USER="tester" \
	-e POSTGRES_PASSWORD="testerpassword" \
	-e POSTGRES_DB="testing_db" \
	postgres:13.1-alpine
	docker run -d --rm --name $(DOCKER_DB)-connector -e PGPASSWORD="testerpassword" postgres \
	psql -U tester -h localhost -d "testing_db"

integration_test: int_db
	@echo "running test files"
	$(GO) test -tags integration $(GO_TEST_FLAGS)

clean:
	docker kill $(DOCKER_DB) && \
	docker rm $(DOCKER_DB)