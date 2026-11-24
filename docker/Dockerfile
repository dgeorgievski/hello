FROM golang:1.27.4@sha256:87f49b77a7b6a99ad7dc656be010ae453322459dd2ed902ac177649a302bc101 AS builder

ARG VERSION
ARG COMMIT
ARG BUILD_DATE
ENV CGO_ENABLED=0
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# The shared CI workflow checks out tags but does not pass build arguments.
# Read Git metadata from the bind mount without copying .git into an image layer.
RUN --mount=type=bind,target=/src \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    set -eu; \
    if [ -z "${VERSION:-}" ] || [ -z "${COMMIT:-}" ] || [ -z "${BUILD_DATE:-}" ]; then \
        if git -c safe.directory=/src rev-parse --verify 'HEAD^{commit}' >/dev/null 2>&1; then \
            tag=$(git -c safe.directory=/src describe --tags --exact-match --match 'v[0-9]*.[0-9]*.[0-9]*' HEAD 2>/dev/null || true); \
            VERSION=${VERSION:-${tag#v}}; \
            COMMIT=${COMMIT:-$(git -c safe.directory=/src rev-parse --short HEAD)}; \
            BUILD_DATE=${BUILD_DATE:-$(git -c safe.directory=/src show -s --format=%cI HEAD)}; \
        fi; \
    fi; \
    module=$(go list -m); \
    go build -buildvcs=false \
        -ldflags "-s -w -X $module/pkg/version.VERSION=${VERSION:-dev} -X $module/pkg/version.COMMIT=${COMMIT:-none} -X $module/pkg/version.DATE=${BUILD_DATE:-none}" \
        -o /hello .

FROM scratch

COPY --from=builder /hello /hello

EXPOSE 8080

CMD ["/hello", "server" ]
