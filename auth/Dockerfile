# Dockerfile для auth сервиса
FROM golang:1.22.4

WORKDIR /auth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth ./cmd/main.go

COPY migrations /auth/migrations

CMD [ "./auth" ]