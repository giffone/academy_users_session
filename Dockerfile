FROM golang:1.21.0-alpine

WORKDIR /web

COPY . .

RUN go mod download
RUN go build -o session_manager cmd/session_manager/main.go

EXPOSE 8080

CMD ["./session_manager"]