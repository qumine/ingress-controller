FROM ubuntu as prep

FROM scratch

EXPOSE 80
EXPOSE 25565

COPY  --from=prep /etc/passwd /etc/passwd
COPY  qumine-ingress qumine-ingress

USER nobody
ENTRYPOINT [ "./qumine-ingress" ]