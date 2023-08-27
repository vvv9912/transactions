BINARY_NAME = transactions
build:
	@echo "  >  Building binary..."
	go build -o ${BINARY_NAME} main.go
run:
	@echo "  >  Run..."
	@echo "  >  Request examples: "

	./${BINARY_NAME}
clean:
	@echo "  >  Cleaning build cache"
	go clean
	rm ${BINARY_NAME}