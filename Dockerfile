FROM golang:1.13.6-alpine

WORKDIR /go
COPY . /go

RUN apk update && apk add git && \
	go get github.com/labstack/echo && \
	go get github.com/labstack/echo/middleware && \
	go get github.com/go-sql-driver/mysql && \
	go get github.com/jinzhu/gorm && \
	go get github.com/joho/godotenv

CMD ["go", "run", "main.go"]