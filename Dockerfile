FROM golang:1.13 AS stage
ARG PACKAGE_NAME
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-w -s" -o /multi-tag-searcher -a -installsuffix cgo cmd/*.go

FROM alpine:3.9
COPY /stortags /stortags
COPY --from=stage /multi-tag-searcher /opt/application
ENTRYPOINT [ "/opt/application" ]
