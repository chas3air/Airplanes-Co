FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o management_flights ./cmd/app

EXPOSE 12002

CMD ["./management_flights"]