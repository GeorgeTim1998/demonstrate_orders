FROM golang:1.21.5

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Add docker-compose-wait tool -------------------
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

RUN go build -o main .

EXPOSE 8080
