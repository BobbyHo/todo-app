FROM golang:1.14.3
RUN mkdir -p /todoapp
WORKDIR /todoapp
COPY todoserver .
COPY todoserverconf.json .
ENTRYPOINT ["./todoserver"]
CMD ["todoserverconf.json"]