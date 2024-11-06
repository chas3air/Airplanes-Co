FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o order_ticket ./cmd/app

EXPOSE 12009

CMD ["./order_ticket"]