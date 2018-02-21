APP_NAME=gtr

all:
	go build -o ${APP_NAME} .
race-build:
	go build -o ${APP_NAME} -race .
