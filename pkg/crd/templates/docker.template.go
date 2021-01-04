package template

const DockerfileTemp = `FROM golang:1.15.6-alpine3.12 as builder

ENV GOPROXY=https://goproxy.cn

ENV GO111MODULE=on

WORKDIR /workspace

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . ./

RUN go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-s -w" -o controller cmd/main.go && chmod +x controller

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM hextechs/alpine:3.12.1

USER root

COPY --from=builder /workspace/controller .

ENTRYPOINT [ "./controller" ]`
