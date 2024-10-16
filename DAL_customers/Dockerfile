FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dal_customers ./cmd/app

EXPOSE 8057

CMD ["./dal_customers"]