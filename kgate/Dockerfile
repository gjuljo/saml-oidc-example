FROM golang:1.12.5-alpine as builder
RUN apk update && apk add --no-cache git build-base make gcc ca-certificates tzdata sqlite && update-ca-certificates

RUN adduser -D -g '' appuser

# Download dependencies
RUN git clone https://github.com/keycloak/keycloak-gatekeeper.git -b 7.0.1 /gatekeeper

WORKDIR /gatekeeper

RUN touch /tmp/bolt

# Copy the source code
RUN go mod download
COPY *.go ./

# And compile the project
RUN make static

FROM scratch
WORKDIR /app

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser /tmp/bolt /tmp/bolt

# Copy our static executable
COPY --from=builder /gatekeeper/bin/keycloak-gatekeeper /app/keycloak-gatekeeper

# Use an unprivileged user.
USER appuser

# Run the binary
ENTRYPOINT ["./keycloak-gatekeeper"]

