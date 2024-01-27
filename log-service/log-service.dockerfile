FROM alpine:latest

RUN mkdir /app

COPY logService /app

CMD ["/app/logService"]