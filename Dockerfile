FROM alpine:3.6

COPY temp /myapp/
COPY uploads /myapp/
COPY myapp /myapp/

CMD ["/myapp/myapp"]