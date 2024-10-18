FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dal_flights ./cmd/app

EXPOSE 12001

CMD ["./dal_flights"]