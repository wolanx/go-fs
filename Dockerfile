FROM golang:1.17.5 AS builder

#ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app

#ADD go.mod .
#ADD go.sum .
#RUN go mod download
COPY . .
#RUN go mod tidy
RUN CGO_ENABLED=0 go build -o go-fs

FROM alpine

LABEL author=github.com/wolanx
ENV TZ utc-8

WORKDIR /app

COPY --from=builder /app/go-fs .

EXPOSE 8080

CMD ["./go-fs"]

# docker build -f Dockerfile -t wolanx/test .
# docker run --restart=unless-stopped --name wt -d -p 8080:8080 wolanx/test
