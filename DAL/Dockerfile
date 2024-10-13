FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dal ./cmd/app

EXPOSE 8056

CMD ["./dal"]