# Start from golang base image
FROM golang:1.19-alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy and download dependency using go mod.
COPY go.mod go.sum ./

RUN go mod download

# Copy the code into the container.
COPY . .

# Build the Go app
RUN go build -o ./reminder ./cmd/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /opt/

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder /app/reminder .
COPY --from=builder /app/config.yaml .

EXPOSE 8000

CMD /opt/reminder