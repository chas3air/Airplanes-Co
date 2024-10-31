FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o purchased_tickets ./cmd/app

EXPOSE 12010

CMD ["./purchased_tickets"]