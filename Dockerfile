FROM ubuntu:latest
LABEL authors="zhelagin.egor"

ENTRYPOINT ["top", "-b"]