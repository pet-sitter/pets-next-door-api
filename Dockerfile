# Build stage
FROM golang:1.20.5 AS build-stage
WORKDIR /app
COPY go.mod *go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server/main.go

# Test stage
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Release stage
FROM alpine:3.17.3 AS build-release-stage
WORKDIR /
COPY --from=build-stage /app/server /app
EXPOSE 8080
RUN adduser -D nonroot
USER nonroot:nonroot
ENTRYPOINT ["/app"]