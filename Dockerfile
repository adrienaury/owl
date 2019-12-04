FROM golang:1.13 AS builder

ENV GOFLAGS="-mod=readonly"

RUN mkdir /home/owl

RUN mkdir -p /workspace
WORKDIR /workspace

ARG GOPROXY

COPY go.* /workspace/
RUN go mod download

COPY . /workspace

ARG VERSION
ARG BUILD_BY

RUN make release

FROM gcr.io/distroless/base

COPY --from=builder /workspace/bin/* /
COPY --from=builder /home/owl /home/owl

WORKDIR /home/owl

ENTRYPOINT [ "/owl" ]

# Build arguments
ARG IMAGE_NAME
ARG IMAGE_TAG
ARG IMAGE_REVISION
ARG IMAGE_DATE

# OCI labels (https://github.com/opencontainers/image-spec/blob/master/annotations.md)
LABEL org.opencontainers.image.created="${IMAGE_DATE}"
LABEL org.opencontainers.image.authors="Adrien AURY <adrien.aury@cgi.com>"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/adrienaury/owl"
LABEL org.opencontainers.image.documentation="https://github.com/adrienaury/owl"
LABEL org.opencontainers.image.source="https://github.com/adrienaury/owl"
LABEL org.opencontainers.image.version="${IMAGE_TAG}"
LABEL org.opencontainers.image.revision="${IMAGE_REVISION}"
LABEL org.opencontainers.image.vendor="Adrien AURY"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.ref.name="${IMAGE_NAME}:${IMAGE_TAG}"
LABEL org.opencontainers.image.title="${IMAGE_NAME}"
LABEL org.opencontainers.image.description="Owl is a set of tools to manage realms of units, users and groups."
