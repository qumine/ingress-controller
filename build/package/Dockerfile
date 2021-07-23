FROM alpine as passwd

FROM scratch

EXPOSE 80
EXPOSE 25565

COPY  ingress-controller ingress-controller
COPY --from=passwd /etc/passwd /etc/passwd
USER nobody

ENTRYPOINT [ "./ingress-controller" ]