FROM alpine:latest
ENV WORKDIR /work
WORKDIR $WORKDIR
COPY gitee-hook ./
EXPOSE 13000
ENTRYPOINT ["sh", "-c", "./gitee-hook -f=/etc/config.yml"]