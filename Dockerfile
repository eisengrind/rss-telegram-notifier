FROM golang:stretch
ENV SRC_REPO=https://github.com/Eisengrind/rss-telegram-notifier.git
ENV SRC=github.com/Eisengrind/rss-telegram-notifier

USER root
RUN cd $GOPATH/src && git clone ${SRC_REPO} ${SRC}
RUN go get ${SRC}/...
RUN cd $GOPATH/src/${SRC} && go install
RUN rm -Rf $GOPATH/src/*

WORKDIR /

ENTRYPOINT [ "rss-telegram-notifier" ]
