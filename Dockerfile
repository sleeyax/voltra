FROM golang:1.22-alpine as builder

RUN apk add --no-cache gcc g++

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# cgo must be enabled for sqlite3 to work
RUN CGO_ENABLED=1 go build -ldflags "-s -w" -o voltra ./cmd/main.go

FROM golang:1.22-alpine

WORKDIR /bot

COPY --from=builder /build/voltra .
COPY LICENSE .
COPY README.md .

CMD ["/bot/voltra"]