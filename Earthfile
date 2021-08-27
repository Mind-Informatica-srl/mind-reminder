all:
    BUILD +lint
   
buildEnvironment:
    FROM golang:1.16.3-alpine
    WORKDIR /work

    # Golangci-lint installation
    RUN apk add curl git
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.39.0

    # Set necessary environmet variables needed for our image
    ENV GO111MODULE=on 
    ENV CGO_ENABLED=0 
    ENV GOOS=linux 
    ENV GOARCH=amd64 
    ENV GOPRIVATE=bitbucket.org/leonardoce

    # Download dependencies
    COPY ./go.mod ./go.sum .
    RUN go mod download

    # Copy source code
    COPY ./pkg pkg
    COPY ./internal internal

lint:
    FROM +buildEnvironment
    ENV CGO_ENABLED=0
    RUN ./bin/golangci-lint run --timeout 5m0s ./...
