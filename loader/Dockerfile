# Build the worker
FROM golang:1.19 AS build-stage

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /loader

# Run the worker
FROM alpine AS run-stage

WORKDIR /run

COPY --from=build-stage /loader /loader

ENTRYPOINT ["/loader"]
