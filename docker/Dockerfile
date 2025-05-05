FROM golang:1.24@sha256:30baaea08c5d1e858329c50f29fe381e9b7d7bced11a0f5f1f69a1504cdfbf5e AS builder

ARG GITHUB_REPO=github.com/altimetrik-digital-enablement-demo-hub/dgeorgievski-hello
ARG VERSION="dev"
ARG COMMIT="none"
ARG BUILD_UNIX_TIME=$(date +%s%M)
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

WORKDIR /app

COPY . .

RUN --mount=type=cache,target=/gomod-cache\
     go mod download

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/gomod-cache \ 
    CGO_ENABLED=0 \
    go build \
    -ldflags " -s -w \
        -X ${GITHUB_REPO}/pkg/version.VERSION=${VERSION} \
        -X ${GITHUB_REPO}/pkg/version.COMMIT=${COMMIT} \
        -X ${GITHUB_REPO}/pkg/version.DATE=${BUILD_UNIX_TIME}" \
    -tags netgo,osusergo -a -o bin/hello

FROM scratch

COPY --from=builder /app/bin/hello /hello

EXPOSE 8080

CMD ["/hello", "server" ]
