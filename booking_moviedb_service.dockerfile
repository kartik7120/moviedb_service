# syntax=docker/dockerfile:1.4
FROM alpine:3.19
WORKDIR /app

COPY moviedbApp /app/moviedbApp
RUN adduser -D appuser && chmod +x /app/moviedbApp

ARG DSN
ARG ENV
ARG STRAPI_API_TOKEN
ARG STRAPI_URL
ARG RABBITMQ_URL
ENV DSN=$DSN
ENV STRAPI_API_TOKEN=$STRAPI_API_TOKEN
ENV STRAPI_URL=$STRAPI_URL
ENV ENV=$ENV
ENV RABBITMQ_URL=$RABBITMQ_URL

# Debug: check secrets

USER appuser

EXPOSE 1102

ENTRYPOINT ["./moviedbApp"]
