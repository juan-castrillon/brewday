FROM golang:1.24-bookworm as builder
ARG VERSION
ENV GO_PROXY=https://proxy.golang.org
WORKDIR /app
COPY . .
RUN VERSION=${VERSION} BUILD_OUT=./brewday make build

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/brewday /brewday
EXPOSE 8080
ENTRYPOINT ["/brewday"]