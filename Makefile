BUILD_DIR=build

all: docker

.PHONY: build
build:
	@ CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BUILD_DIR}/app ./example/main.go

docker: stop
	@ echo "Starting app in docker..."
	@ docker-compose build --no-cache && docker-compose up -d

stop:
	@ echo "Stopping app..."
	@ docker-compose down -v --remove-orphans # --rmi=all

test:
	@ go test -v ./...
