FROM golang:1.20 as build
WORKDIR /app
ARG SRC=/cmd
RUN --mount=source=go.mod,target=./go.mod \
    #--mount=source=go.sum,target=./go.sum \
    --mount=type=cache,target=$GOPATH/pkg/mod \
    go mod download && go mod verify

RUN --mount=source=$SRC,target=$SRC \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /app/bin/app $SRC

FROM gcr.io/distroless/base-debian11 as release-debian
WORKDIR /
COPY --from=build /app/bin/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]