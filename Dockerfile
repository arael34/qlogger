# !! This isn't finished

FROM golang:latest

# Create working directory
WORKDIR /app

# Copy dependency management into dir
COPY go.mod go.sum ./

# Download deps
RUN go mod download

# Copy root dir into working dir
COPY . /app/.

# Build app
RUN go build -o out ./cmd/

# Run app
CMD ["./out"]
