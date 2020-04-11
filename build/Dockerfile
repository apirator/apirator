## Builder Image
FROM golang:1.14-stretch as builder
ENV GO111MODULE=on
COPY . /apirator
WORKDIR /apirator
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o build/_output/bin/apirator cmd/manager/main.go

## Run Image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/apirator \
    USER_UID=1001 \
    USER_NAME=apirator

# install operator binary
COPY --from=builder /apirator/build/_output/bin/apirator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
