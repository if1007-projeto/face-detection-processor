FROM golang:1.12

ARG APP_NAME=face-detection-processor
ARG APP_DIR=/opt/${APP_NAME}
ARG APP=${APP_DIR}/${APP_NAME}

WORKDIR ${APP_DIR}

COPY *.go ./

RUN go get -u github.com/esimov/pigo/cmd/pigo && \
    go get -u github.com/fogleman/gg && \
    go get -u github.com/disintegration/imaging && \
    go get -u github.com/lovoo/goka && \
    go build

CMD [ "${APP}" ]
