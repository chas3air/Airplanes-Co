# Set the Go version as an argument
ARG GO_VERSION=latest
FROM golang:${GO_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/purchased_tickets ./cmd/app

FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates tzdata

COPY --from=build /bin/purchased_tickets /bin/

RUN adduser -D appuser
USER appuser

EXPOSE 12010

ENTRYPOINT ["/bin/purchased_tickets"]