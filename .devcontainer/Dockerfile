# Use the official Go image with the latest version
FROM golang:latest

# Set the working directory inside the container
WORKDIR /akhanda

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application's source code to the container
COPY . .

# Install any necessary Go tools (optional, customize as needed)
RUN go install golang.org/x/tools/gopls@latest

