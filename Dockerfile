FROM golang:1.22.8@sha256:0ca97f4ab335f4b284a5b8190980c7cdc21d320d529f2b643e8a8733a69bfb6b
RUN ls /opt
COPY . /opt/
WORKDIR /opt
RUN make install
RUN make cmd/web/build
CMD ["make", "cmd/web/run"]
