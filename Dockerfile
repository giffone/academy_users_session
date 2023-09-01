FROM golang:1.21.0-alpine

WORKDIR /web

COPY . .

RUN mod download
RUN go build -o session_manager cmd/session_manager/main.go

CMD ["./session_manager"]