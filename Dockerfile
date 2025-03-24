# setup project and deps
FROM golang:1.24-bookworm AS init

WORKDIR /go/go-find-liquor/

COPY go.mod* go.sum* ./
RUN go mod download

COPY . ./

FROM init AS vet
RUN go vet ./...

# run tests
FROM init AS test
RUN go test -coverprofile c.out -v ./...

# build binary
FROM init AS build
ARG LDFLAGS

RUN CGO_ENABLED=0 go build -ldflags="${LDFLAGS}"

# runtime image including CA certs and tzdata
FROM gcr.io/distroless/static-debian12:latest
# Copy our static executable.
COPY --from=build /go/go-find-liquor/go-find-liquor /go/bin/go-find-liquor
# Expose port for publishing as web service
# EXPOSE 8081
# Run the binary.
ENTRYPOINT ["/go/bin/go-find-liquor", "--config", "/config.yaml"]
