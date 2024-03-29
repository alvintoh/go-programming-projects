# Start from the official golang image
FROM golang:1.21.3 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set a build argument for the service name
ARG SERVICE_NAME

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/${SERVICE_NAME}/main.go

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the application files from the resource directory
COPY --from=builder /app/resource/* ./resource/

# Copy the Swagger files from the swagger directory
COPY --from=builder /app/swagger ./swagger/

# Set CONFIG_ENV environment variables
ENV CONFIG_ENV=application-deploy

# Command to run the executable
CMD ["./main"]