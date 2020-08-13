FROM alpine as passwed

FROM scratch

EXPOSE 80
EXPOSE 25565

COPY  qumine-ingress qumine-ingress
COPY --from=passwed /etc/passwd /etc/passwd
USER nobody

ENTRYPOINT [ "./qumine-ingress" ]