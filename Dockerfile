ARG GO_VERSION=1.17
ARG APP_NAME="knowledge"
ARG PORT=8000

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY ./server/go.mod ./server/go.sum ./
RUN go mod download

COPY ./server ./

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o /app .

FROM scratch

WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
EXPOSE ${PORT}

ENTRYPOINT ["/app"]