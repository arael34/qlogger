FROM golang:latest

# Set cwd
WORKDIR /app/

# Copy everything from the current directory to the working directory
COPY . /app/.

# Download dependencies
RUN go mod download

# Build
RUN go build -o out

# This isn't needed for certain hosting services.
# EXPOSE 8080

# Run
CMD ["./out"]
