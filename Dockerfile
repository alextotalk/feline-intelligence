FROM golang:1.24.1-alpine AS builder

WORKDIR /usr/local/src

 RUN apk --no-cache add bash git make

 COPY go.mod go.sum ./
RUN go mod download

 RUN mkdir -p /usr/local/src/config

 COPY config/local.yaml /usr/local/src/config/local.yaml

 COPY cmd/ cmd/
COPY internal/ internal/
COPY docs/ docs/


RUN go build -o ./bin/app ./cmd/app/main.go


FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app

COPY --from=builder /usr/local/src/config /usr/local/src/config

WORKDIR /usr/local/src

CMD ["/app"]
