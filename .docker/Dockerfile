FROM alpine:latest

WORKDIR /app
VOLUME [ "/data" ]

COPY  ./bin/quotebot  .
RUN ls -lha

ENTRYPOINT [ "./quotebot", "-config", "/data/config.json", "-token", "/data/token", "-dataPath", "/data/store"]