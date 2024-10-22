# syntax=docker/dockerfile:1

################################################################################
ARG GO_VERSION=1.23.0

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

ARG PROJECT_PATH=apps/server
ARG LIBVIPS_VERSION=8.15.5

LABEL org.opencontainers.image.source="https://github.com/taskemapp/server"

RUN --mount=type=bind,source=${PROJECT_PATH}/go.sum,target=go.sum \
    --mount=type=bind,source=${PROJECT_PATH}/go.mod,target=go.mod \
    go mod download -x

RUN DEBIAN_FRONTEND=noninteractive \
    apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libglib2.0-dev \
    libexpat1-dev \
    libjpeg62-turbo-dev \
    libtiff5-dev \
    libarchive-dev \
    meson \
    ninja-build && \
    rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.xz && \
    tar -xf vips-${LIBVIPS_VERSION}.tar.xz && \
    cd vips-${LIBVIPS_VERSION} && \
    meson setup build --prefix /usr/local --libdir lib && \
    meson compile -C build && \
    meson install -C build && \
    ldconfig && \
    cd .. && rm -rf vips-${LIBVIPS_VERSION}*

ARG TARGETARCH

RUN --mount=type=bind,source=./${PROJECT_PATH}/,target=. \
    CGO_ENABLED=1 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server

################################################################################
FROM gcr.io/distroless/base-debian12:nonroot AS final

COPY --from=build /usr/local/lib /usr/local/lib

ENV LD_LIBRARY_PATH="/usr/local/lib"

COPY --from=build /bin/server /bin/

COPY migrations migrations

EXPOSE 50051

ENTRYPOINT [ "/bin/server" ]
