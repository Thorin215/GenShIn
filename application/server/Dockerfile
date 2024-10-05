FROM golang:1.18 AS server
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/genshin
COPY . .
RUN CGO_ENABLED=0 go build -v -o "server" .

FROM scratch
WORKDIR /root/genshin/
COPY --from=server /root/genshin/server ./
COPY config.yaml ./
ENTRYPOINT ["./server"]