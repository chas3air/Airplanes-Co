FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o management_tickets ./cmd/app

EXPOSE 12008

CMD ["./management_tickets"]