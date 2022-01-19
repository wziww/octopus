FROM golang:1.14.3-alpine3.11

ENV APP_DIR $GOPATH/src/octopus
RUN apk add busybox-extras curl make && mkdir -p $APP_DIR && mkdir -p /usr/local/octopus/data/
ADD . $APP_DIR
WORKDIR $APP_DIR 

RUN make build
EXPOSE 8089
CMD ["$@"]
