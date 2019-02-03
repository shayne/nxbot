############################
# STEP 1 build executable binary
############################
FROM golang@sha256:9657ef82d7ead12e0c88c7f4708e78b50c5fd3c1893ac0f2f0924ab98873aad8 as builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
# Create nxbot
RUN adduser -D -g '' nxbot
WORKDIR /nx-code
COPY . .
# Fetch dependencies.
# Using go mod with go 1.11
RUN go mod download
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/nxbot github.com/jacknx/nxbot/cmd/nxbot
############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/nxbot /go/bin/nxbot
# Use an unprivileged user.
USER nxbot
# Run the nxbot binary.
ENTRYPOINT ["/go/bin/nxbot"]