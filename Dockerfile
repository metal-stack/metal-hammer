FROM golang:1.11-stretch as builder

RUN apt update && apt -y install make git

WORKDIR /app
COPY go.mod .
RUN go mod download
COPY .git ./.git
COPY cmd ./cmd
COPY main.go ./main.go
COPY Makefile ./Makefile
RUN make

FROM alpine:3.8
COPY --from=builder /app/bin/discover /bin/discover
RUN apk update \
 && apk add \
        ca-certificates \
        lshw \
        sgdisk \
        parted \
 && wget https://github.com/genuinetools/img/releases/download/v0.5.0/img-linux-amd64 -O /bin/img \
 && chmod +x /bin/img \
 && wget https://github.com/opencontainers/runc/releases/download/v1.0.0-rc5/runc.amd64 -O /bin/runc \
 && chmod +x /bin/runc

CMD ["/bin/discover"]
