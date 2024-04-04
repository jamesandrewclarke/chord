FROM golang:1.21 as build

WORKDIR /usr/src/app

RUN apt-get update && apt-get install protobuf-compiler build-essential -y 

# Copy go mod and go sum first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy protos
COPY protos/chord.proto protos/ 
COPY Makefile .
RUN make

# Copy source
COPY chord chord
COPY protos protos
COPY peer peer

# Compile
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./peer

# Final result is a bare container with just the binary, results in a much smaller image
FROM scratch

COPY --from=0 /usr/local/bin/app /app

ENTRYPOINT ["/app"]
CMD ["-port", "8080"]

EXPOSE 8080