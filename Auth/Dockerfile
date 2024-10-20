FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth ./cmd/app

EXPOSE 12005

CMD ["./auth"]