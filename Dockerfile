FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ /app/

COPY .env /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /easy-tasks-api ./cmd/easy-tasks-api/main.go

EXPOSE 3000

CMD ["/easy-tasks-api"]