ARG ARCH=
FROM ${ARCH}alpine

EXPOSE 80
EXPOSE 25565

COPY  ingress-controller ingress-controller
USER nobody

ENTRYPOINT [ "./ingress-controller" ]