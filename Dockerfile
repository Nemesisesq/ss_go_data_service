FROM golang:1.7-alpine
MAINTAINER Carl Lewis "carl@streamsavvy.tv"


WORKDIR /go/src/app

RUN mkdir -p /root/.ssh
ADD docker-repo-key /root/.ssh/docker-repo-key
RUN chmod 700 /root/.ssh/docker-repo-key
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config

COPY . /go/src/app

RUN go get -d -v

RUN go install -v

CMD ["app"]