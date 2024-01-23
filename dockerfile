FROM golang:alpine AS builder

WORKDIR /usr/local/src

ADD go.mod .

COPY . .

RUN go build -o ./bin/auth cmd/auth/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/auth .
COPY config/local.yaml .

ENV CONFIG_PATH=./local.yaml

CMD ["/auth"]