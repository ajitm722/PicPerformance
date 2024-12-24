# Stage 1: Build the application
FROM golang:1.22.1-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /my-app

# Stage 2: Create a minimal runtime image
FROM alpine:3.18
COPY --from=build /my-app /my-app
ENTRYPOINT ["/my-app"]

