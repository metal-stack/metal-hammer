FROM golang:1.11-stretch as builder

RUN apt update \
 && apt -y install make git

WORKDIR /app
COPY go.mod .
RUN go mod download
COPY .git ./.git
COPY cmd ./cmd
COPY main.go ./main.go
COPY Makefile ./Makefile
RUN make

FROM alpine:3.8
LABEL maintainer FI-TS Devops <devops@f-i-ts.de>
COPY --from=builder /app/bin/discover /bin/discover
RUN apk update \
 && apk add \
        ca-certificates \
        lshw \
        sgdisk \
        e2fsprogs

CMD ["/bin/discover"]
