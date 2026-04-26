FROM golang:1.26.1
RUN mkdir -p /todoapp
WORKDIR /todoapp
COPY todoserver .
#COPY todoserverconf.json .
ENTRYPOINT ["./todoserver"]
CMD ["todoserverconf.json"]
