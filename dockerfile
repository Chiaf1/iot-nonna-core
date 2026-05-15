# Stage 1: BUILD
FROM --platform=$BUILDPLATFORM golang:1.26.3 AS builder

# mod caching 
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# copying surce
COPY . .

# cross-compile respecting target platform
ARG TARGETOS
ARG TARGETARCH

# building the app
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o app ./cmd/iot-nonna-core

# Stage 2: final image
FROM  alpine:3.23.3

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY migrations ./migrations
COPY docs ./docs
COPY --from=builder /build/app .

ENV RUN_MIGRATIONS=false
ENV RUN_SEEDING=false

EXPOSE 3030

CMD ["./app"]