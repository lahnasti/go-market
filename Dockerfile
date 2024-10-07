FROM golang:1.22.4

WORKDIR /market

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

COPY migrations /market/migrations

CMD [ "./main" ]