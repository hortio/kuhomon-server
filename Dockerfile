FROM alpine:3.8
LABEL maintainer="Sergei Silnov <po@kumekay.com>"

RUN apk -U add ca-certificates

EXPOSE 8080

COPY kuhomon /bin/kuhomon

CMD ["kuhomon"]