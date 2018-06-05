FROM alpine:3.7

WORKDIR /app
COPY . /app/

EXPOSE 8080

CMD ["./go-fs"]
