FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY moviedbApp /app/moviedbApp

RUN chmod +x moviedbApp

CMD [ "./moviedbApp" ]