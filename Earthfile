all:
    BUILD +lint
   
buildEnvironment:
    FROM golang:1.16.3-alpine
    WORKDIR /work

    # Golangci-lint installation
    RUN apk add curl git
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.39.0

    RUN git config --global url."https://simonerocchi:yq3RGzmhzSKbyM5hAhRF@bitbucket.org/".insteadOf "https://bitbucket.org/"

    # Set necessary environmet variables needed for our image
    ENV GO111MODULE=on 
    ENV CGO_ENABLED=0 
    ENV GOOS=linux 
    ENV GOARCH=amd64 
    ENV GOPRIVATE=bitbucket.org/leonardoce

    # Download dependencies
    COPY ./lamicolor-web-server/go.mod ./lamicolor-web-server/go.sum .
    RUN go mod download

    # Copy source code
    COPY ./lamicolor-web-server/cmd cmd
    COPY ./lamicolor-web-server/internal internal

lint:
    FROM +buildEnvironment
    ENV CGO_ENABLED=0
    RUN ./bin/golangci-lint run --timeout 5m0s ./...
