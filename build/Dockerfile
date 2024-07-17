FROM alpine:latest

RUN mkdir /app
COPY ./dist/icbc-checker /app/

WORKDIR /app
CMD ["./icbc-checker"]