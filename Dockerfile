FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -ldflags "-s -w" -o voltra ./cmd/main.go

FROM golang:1.22-alpine

WORKDIR /bot

COPY --from=builder /build/voltra .
COPY LICENSE .
COPY README.md .

CMD ["/bot/voltra"]