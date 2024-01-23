FROM alpine:latest

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai
ARG ENV
ARG FILENAME
ENV ENV=$ENV
ENV FILENAME=$FILENAME
ENV PATH="../../../"
WORKDIR /app

COPY  ${PATH}configs/ ./
COPY ${PATH}bin/${FILENAME} ./app

CMD ["./app", "-conf", "config.yaml"]


