FROM scratch
ADD config.json /
ADD ./bin/main /
CMD ["/main"]
