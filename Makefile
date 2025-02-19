TEST_DIR := ./tests

test:
	@echo "Running tests in $(TEST_DIR)..."
	go test -v $(TEST_DIR)/...