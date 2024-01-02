# Stage 1: Build the Go application
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/go-rest-api ./

# Stage 2: Use a more efficient base image for KeyDB
FROM alpine:latest as keydb-builder

# Install build dependencies, uuid-dev and libcurl4-openssl-dev
RUN apk --no-cache add \
    build-base \
    tcl \
    git \
    wget \
    uuid-dev \
    libcurl

# Clone and build KeyDB
RUN git clone --branch v6.0.16 https://github.com/EQ-Alpha/KeyDB.git /tmp/keydb \
    && cd /tmp/keydb \
    && make \
    && make install

# Stage 3: Use Alpine for the final image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    fuse \
    libuuid \
    libcurl \
    wget

WORKDIR /root/

# Copy KeyDB binaries
COPY --from=keydb-builder /usr/local/bin/keydb-server /usr/local/bin/

# Create directories for LiteFS and KeyDB within the mounted volume
RUN mkdir -p /var/lib/litefs /var/lib/keydb_storage/data

# Copy the compiled application from the builder stage
COPY --from=builder /app/go-rest-api .
COPY litefs.yml ./

# Copy the KeyDB configuration file
COPY keydb.conf /etc/keydb.conf

# Expose the necessary ports
EXPOSE 8080 8081 6379

# Environment variables for KeyDB
ENV KEYDB_SERVERS=keydb-internal:6379
ENV KEYDB_PASSWORD=yourpassword

# Start KeyDB, mount LiteFS, and run the application
CMD keydb-server /etc/keydb.conf --daemonize yes && litefs mount

# Our final Docker image stage starts here.
FROM alpine
ARG LITEFS_CONFIG=litefs.yml

# Copy binaries from the previous build stages.
COPY --from=flyio/litefs:0.5.8 /usr/local/bin/litefs /usr/local/bin/litefs
COPY --from=builder /usr/local/bin/go-rest-api /usr/local/bin/go-rest-api

# Copy the possible LiteFS configurations.
ADD litefs.yml /tmp/litefs.yml

# Move the appropriate LiteFS config file to /etc/
RUN cp /tmp/$LITEFS_CONFIG /etc/litefs.yml

# Setup our environment to include FUSE & SQLite
RUN apk --no-cache add bash fuse3 sqlite ca-certificates curl

# Run LiteFS as the entrypoint
ENTRYPOINT ["litefs", "mount"]
