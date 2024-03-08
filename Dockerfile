FROM alpine
COPY echo-server /
ENTRYPOINT ["/echo-server"]