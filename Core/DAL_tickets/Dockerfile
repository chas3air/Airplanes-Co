FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dal_tickets ./cmd/app

EXPOSE 12002

CMD ["./dal_tickets"]