FROM alpine:3.6

WORKDIR /app
COPY . /app/

CMD ["./app"]
