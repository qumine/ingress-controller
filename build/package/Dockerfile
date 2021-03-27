FROM alpine as passwed

FROM scratch

EXPOSE 80
EXPOSE 25565

COPY  ingress-controller ingress-controller
COPY --from=passwed /etc/passwd /etc/passwd
USER nobody

ENTRYPOINT [ "./ingress-controller" ]