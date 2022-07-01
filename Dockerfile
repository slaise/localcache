# Base image, golang alpine(1.18.3)
FROM golang:alpine as builder
WORKDIR /workspace
# Copy all files into the image
COPY go.mod go.mod
COPY go.sum go.sum
COPY localcache/ localcache/
COPY main.go main.go
# Run go mod
RUN go mod download
# Build Go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go

FROM debian AS runner
WORKDIR /go/localcache
# COPY executable file from previous builder stage
COPY --from=builder /workspace/main .
# Expose ports
EXPOSE 8000
# Run Go program, just like locally
ENTRYPOINT ["./main"]
