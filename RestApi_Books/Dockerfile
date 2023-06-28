# FROM golang:latest

# WORKDIR /app
# COPY go.mod .
# COPY go.sum .


# RUN go mod download
# COPY . .


# EXPOSE 8080
# RUN go build
# CMD [ "./modules" ]

# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./

# Expose the port for the API
EXPOSE 8080

# Set the entry point of the container
CMD ["./main"]