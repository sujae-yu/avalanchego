# The version is supplied as a build argument rather than hard-coded
# to minimize the cost of version changes.
ARG GO_VERSION=INVALID # This value isn't intended to be used but silences a warning

# AVALANCHEGO_NODE_IMAGE needs to identify an existing node image and should include the tag
ARG AVALANCHEGO_NODE_IMAGE=INVALID # This value isn't intended to be used but silences a warning

# ============= Compilation Stage ================
FROM golang:$GO_VERSION-bullseye AS builder

WORKDIR /build

# Copy and download avalanche dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build xsvm
RUN ./scripts/build_xsvm.sh

# ============= Cleanup Stage ================
FROM $AVALANCHEGO_NODE_IMAGE AS execution

# Copy the xsvm binary to the default plugin dir for images
COPY --from=builder /build/build/xsvm /avalanchego/build/plugins/v3m4wPxaHpvGr8qfMeyK6PRW3idZrPHmYcMTt7oXdK47yurVH

# The node image's entrypoint will be reused.
