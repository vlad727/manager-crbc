FROM alpinelinux/golang AS builder
WORKDIR /app
COPY . /app
USER root
RUN env GOOS=linux GOARCH=amd64 && go build -o manager-crbc /app/cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/manager-crbc /app/
# create dir for html files
RUN mkdir /app/tmpl
# copy html files from tmpl to /app/tmpl
COPY --from=builder /app/tmpl /app/tmpl
RUN apk update --no-check-certificate \
    && apk add --no-check-certificate curl net-tools
RUN ls /app/*
RUN  chmod u+x manager-crbc && mkdir /certs  /files
CMD ["./manager-crbc"]