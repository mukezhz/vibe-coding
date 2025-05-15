FROM golang:1.23.0

# Required because go requires gcc to build
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    inotify-tools \
    curl \
    chromium \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN echo $GOPATH
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN curl -sSf https://atlasgo.sh | sh

COPY . /clean_web
WORKDIR /clean_web
RUN go mod download

CMD sh /clean_web/docker/run.sh
