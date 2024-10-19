FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o flights_catalog ./cmd/app

EXPOSE 12001

CMD ["./flights_catalog"]