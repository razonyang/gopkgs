FROM golang:latest as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY ./ /src

RUN go mod download

RUN go build github.com/gobuffalo/packr/v2/packr2

RUN ./packr2 build -o dist/main


FROM ubuntu

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

COPY --from=builder /src/dist/main /app/main

# Command to run
ENTRYPOINT ["/app/main", "serve"]