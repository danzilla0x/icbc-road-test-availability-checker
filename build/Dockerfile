FROM alpine:latest

RUN mkdir /app
COPY config.env ./dist/icbc-checker /app/

WORKDIR /app
CMD ["./icbc-checker"]