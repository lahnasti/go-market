# Dockerfile для products сервиса
FROM golang:1.23.1

WORKDIR /market

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o market ./cmd/main.go

COPY migrations /market/migrations

CMD [ "./market" ]