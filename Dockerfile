FROM alpine:3.6

WORKDIR /myapp
COPY . /myapp/

CMD ["./myapp"]

