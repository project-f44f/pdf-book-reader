FROM m.daocloud.io/docker.io/golang:1.22.9 AS builder

ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod ./
COPY main.go ./
RUN go build -o server

FROM m.daocloud.io/docker.io/debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/server .

COPY pdfjs-5.1.91-dist/web ./web

COPY pdfjs-5.1.91-dist/build ./build

RUN mkdir /app/pdfs

EXPOSE 8000

CMD ["./server"]
