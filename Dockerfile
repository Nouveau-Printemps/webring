FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY . .

RUN bun i && bun run build

FROM golang:1.24-alpine

WORKDIR /app

COPY --from=builder /app .

RUN go mod tidy && go build -o app .

EXPOSE 80

ENV CONFIG="/config.toml"

CMD ./app -config $CONFIG