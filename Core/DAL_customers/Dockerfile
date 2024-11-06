FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dal_customers ./cmd/app

EXPOSE 12000

CMD ["./dal_customers"]