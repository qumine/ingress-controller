FROM alpine

EXPOSE 80
EXPOSE 25565

COPY  qumine-ingress qumine-ingress

USER nobody
ENTRYPOINT [ "./qumine-ingress" ]