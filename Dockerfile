# Use the official Go image as a builder
FROM golang:1.20 as builder
# Set the working directory in the builder container
WORKDIR /app
# Copy the entire project into the workdir
COPY . .
# Build the Go app. Since your main.go is inside the cmd folder,
# we need to adjust the path.
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/
# Use a lightweight Alpine image for the runtime container
FROM alpine:latest
# Add necessary certificates for secure communication
RUN apk --no-cache add ca-certificates
# Set the working directory in the container
WORKDIR /root/
# Copy the binary built in the builder stage into the runtime container
COPY --from=builder /app/main .
# Command to run the application
CMD ["./main"]