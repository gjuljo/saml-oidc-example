FROM golang:1.12.5-alpine AS gobuilder
RUN apk update && apk add --no-cache git build-base make gcc ca-certificates tzdata && update-ca-certificates

RUN adduser -D -g '' appuser

WORKDIR /gosaml/

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./
COPY *.cert ./
COPY *.key ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo      \
       -ldflags="-w -s -X main.BUILDDATE=`date +%Y-%m-%dT%T%z`"  \    
       -o /gosaml/gosamltest

FROM node:12-slim AS elmbuilder

WORKDIR /gosaml-elm/

RUN apt-get update; \
   if [ "$http_proxy" -o "$https_proxy" ]; then \
      wget https://github.com/elm/compiler/releases/download/0.19.0/binary-for-linux-64-bit.gz -P /tmp/ ; \
      gunzip /tmp/binary-for-linux-64-bit.gz -c > /usr/local/bin/elm; \
      chmod ugo+x /usr/local/bin/elm; \
   else npm install -g elm -unsafe-perm=true --allow-root && npm install uglify-js -g ; \
   fi && \
   npm install uglify-js -g

COPY public public
COPY elm.json .
COPY elmsrc elmsrc

# BUILD ELM SOURCE CODE
RUN mkdir -p public/js && elm make elmsrc/Main.elm --output=elm.js \
   && uglifyjs elm.js --compress "pure_funcs=[F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9],pure_getters,keep_fargs=false,unsafe_comps,unsafe" | uglifyjs --mangle --output=public/js/elm.min.js

FROM scratch
EXPOSE 8000
WORKDIR /app

# Import from builder.
COPY --from=gobuilder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=elmbuilder /gosaml-elm/public /app/public
COPY --from=gobuilder /gosaml/gosamltest /app/gosamltest
COPY --from=gobuilder --chown=appuser:0 /gosaml/myservice.key /app/certs/myservice.key
COPY --from=gobuilder --chown=appuser:0 /gosaml/myservice.cert /app/certs/myservice.cert

# Use an unprivileged user.
USER appuser

# Run the hello binary.
ENV CERTFILE /app/certs/myservice.cert
ENV KEYFILE /app/certs/myservice.key

CMD ["./gosamltest"]