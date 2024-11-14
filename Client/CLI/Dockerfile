# Установите версию Go как аргумент
ARG GO_VERSION=latest
FROM golang:${GO_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} go build -o /src/cli ./cmd/app

FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates tzdata

COPY --from=build /src/cli /cli

RUN adduser -D appuser
USER appuser

WORKDIR /

EXPOSE 12016

CMD [ "./cli" ]