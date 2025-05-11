FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY moviedbApp /app/moviedbApp

CMD [ "./moviedbApp" ]