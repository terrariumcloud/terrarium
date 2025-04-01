FROM golang:1.22 AS test
ENV GOPRIVATE=github.com/terrariumcloud
WORKDIR /workspace
COPY . /workspace
RUN go test -v -cover ./...

FROM golang:1.22 AS build
ENV CGO_ENABLED=0 GOOS=linux GARCH=amd64
WORKDIR /workspace
COPY . /workspace
RUN go mod vendor
ARG BUILD_INFO_VERSION
RUN go build  -ldflags "-X 'main.buildInformationVersion=$BUILD_INFO_VERSION'" -o terrarium
RUN apt-get update && \
    apt-get install -y ca-certificates

FROM scratch
COPY --from=build /workspace/terrarium /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT [ "/terrarium" ]
