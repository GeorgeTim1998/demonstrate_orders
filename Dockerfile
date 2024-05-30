FROM golang:1.21.5

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# RUN ./bin/goose -dir db/migrations postgres $DATABASE_URL up
# RUN go build -o main .

# EXPOSE 8080

# CMD ["./main"]

# RUN go run /tests/main_test.go
CMD ["sh", "-c", "while true; do sleep 1; done"]