BINARY_NAME=IMAGE-SERVER

dev: 
	clear
	go run cmd/main.go

build:
	go build -o ${BINARY_NAME} cmd/main.go

run:
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}	