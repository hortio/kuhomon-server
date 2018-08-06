FROM alpine:3.8

RUN apk -U add ca-certificates

EXPOSE 8080

ADD kuhomon /bin/kuhomon

CMD ["kuhomon"]