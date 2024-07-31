FROM alpinelinux/golang AS builder
WORKDIR /app
COPY . /app
USER root
RUN env GOOS=linux GOARCH=amd64 && go build -o manager-crbc /app/cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/ /app/
RUN apk add net-tools
RUN ls /app/*
RUN  chmod u+x manager-crbc && mkdir /certs  /files
CMD ["./manager-crbc"]