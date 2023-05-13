
# Step 2: Builder
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o binfile ./cmd/core 

# Step 3: Final
FROM alpine
COPY --from=builder /app/binfile /binfile
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/binfile"]
