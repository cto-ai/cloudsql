############################
# Build container
############################
FROM golang:1.14 AS build

WORKDIR /go/src

ADD . .
RUN go get ./...
RUN go build -ldflags="-s -w" -o /ops/cloudsql

############################
# Final container
############################
FROM registry.cto.ai/official_images/base:2-stretch-slim

ENV CLOUD_SDK_VERSION=274.0.1
ENV PATH /usr/local/bin/google-cloud-sdk/bin:$PATH

RUN apt update && apt install -y python-pip curl

RUN curl -Os https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
  && tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
  && rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
  && mv google-cloud-sdk/ /usr/local/bin \
  && gcloud components install beta --verbosity="error" \
  && cd /usr/local/bin/google-cloud-sdk/bin \
  && gcloud config set core/disable_usage_reporting true \
  && gcloud config set component_manager/disable_update_check true

ENV GOOGLE_APPLICATION_CREDENTIALS="/ops/gcp.json"

COPY --from=build /ops/cloudsql /ops/cloudsql