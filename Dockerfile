FROM alpine
 
WORKDIR /build
COPY ikuai-bypass .
 
CMD ["./ikuai-bypass", "-c", "/etc/ikuai-bypass/config.yml"]