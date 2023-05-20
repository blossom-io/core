ARG GO=golang:1.20-alpine

# Step 1: Modules caching
FROM ${GO} as deps

WORKDIR /modules
COPY go.mod go.sum .
RUN go mod download

# Step 2: Builder
FROM ${GO} as builder
COPY --from=deps /go/pkg /go/pkg
WORKDIR /app
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -o bin -ldflags=-X=main.version=${VERSION} ./cmd/core

# Step 3: Final
FROM alpine
COPY --from=builder /app/bin /bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["bin"]
