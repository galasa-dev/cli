FROM ubuntu:jammy

RUN apt-get update \
  && apt-get install -y curl wget unzip zip

RUN apt-get update \
  && apt-get install -y fontconfig libfreetype6

RUN mkdir -p /workspace
WORKDIR /workspace

