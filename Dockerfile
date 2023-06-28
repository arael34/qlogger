# This dockerfile runs the portal.
FROM golang:latest

# Set cwd
WORKDIR /app/

# Copy everything from the portal directory to the working directory
COPY ./_portal/ /app/.

# Download dependencies
RUN go mod download

# Build
RUN go build -o out

# This isn't needed for certain hosting services.
# EXPOSE 8080

# Run
CMD ["./out"]
