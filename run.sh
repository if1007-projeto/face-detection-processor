#!/bin/bash

PARAMS=""

# check if kafka url is assigned
if [[ ! -z $KAFKA_URL ]]; then
    # echo "Environment variable \$KAFKA_URL is required, but has no value."
    PARAMS=$PARAMS"-b $KAFKA_URL ";
fi

# check if kafka topic is assigned
if [[ ! -z $FRAMES_TOPIC ]]; then
    # echo "Environment variable \$CONSUMER_TOPIC is required, but has no value.";
    PARAMS=$PARAMS"-ct $FRAMES_TOPIC ";
fi

# check if kafka topic is assigned
if [[ ! -z $FACES_TOPIC ]]; then
    # echo "Environment variable \$PRODUCER_TOPIC is required, but has no value."
    PARAMS=$PARAMS"-ft $FACES_TOPIC ";
fi

# check if kafka topic is assigned
if [[ ! -z $MARKS_TOPIC ]]; then
    # echo "Environment variable \$PRODUCER_TOPIC is required, but has no value."
    PARAMS=$PARAMS"-mt $MARKS_TOPIC ";
fi

echo "running \$ $APP $PARAMS"

sh -c "$APP $PARAMS"