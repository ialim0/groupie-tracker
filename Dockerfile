# Use the official Golang image as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod  ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN go build -o app

# Expose the port that the Go application listens on
EXPOSE 8080

# Set a meta description environment variable
ENV META_DESCRIPTION "Groupie-tracker  app."
LABEL maintainer="ialimoud passar & addiedhiou"
LABEL version="1.0"
LABEL description="Groupie Trackers consists on receiving a given API and manipulate the data contained in it, in order to create a site, displaying the information."

# Run the Go application when the container starts
CMD ["./app"]
