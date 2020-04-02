# Build stage
FROM golang:alpine AS builder
ENV CGO_ENABLED 0
WORKDIR /app
ADD . /app
RUN go build -o server .

# Final stage
FROM alpine
ADD ./migrations /migrations
COPY --from=builder /app/server /
ENTRYPOINT ["/server"]
EXPOSE 50051/tcp
EXPOSE 8090/tcp
