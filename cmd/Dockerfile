ARG GO_VERS=1.15.7
FROM golang:$GO_VERS AS build
WORKDIR /go/src/go-todos
ADD .. .
RUN go mod download && \
    go test -v -cover
RUN go build -o /go/bin/todos

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/todos /
CMD ["/todos"]
