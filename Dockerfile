FROM alpine

COPY http-exporter /
ENTRYPOINT [ "/http-exporter" ]
