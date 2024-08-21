FROM golang:1.22

COPY . .

RUN bin/build

ENTRYPOINT [ "out/blog" ]
