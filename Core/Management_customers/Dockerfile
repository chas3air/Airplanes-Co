FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o management_customers ./cmd/app

EXPOSE 12007

CMD ["./management_customers"]