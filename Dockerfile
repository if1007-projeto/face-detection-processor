FROM golang:1.12

ARG APP_NAME=face-detection-processor
ARG APP_DIR=/opt/${APP_NAME}
ENV MY_GO_APP=${APP_DIR}/${APP_NAME}

WORKDIR ${APP_DIR}

COPY *.go ./

RUN apt-get update && apt install gnutls-bin -y

RUN go get -u github.com/esimov/pigo/cmd/pigo && \
    go get -u github.com/fogleman/gg && \
    go get -u github.com/disintegration/imaging && \
    go get -u github.com/lovoo/goka && \
    go build

CMD [ ${MY_GO_APP} ]
