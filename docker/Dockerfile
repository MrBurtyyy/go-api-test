FROM golang:latest as build

WORKDIR /src
COPY ./ /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/build/myapp ./main.go

FROM alpine
COPY --from=build /src/build/myapp /src/bin/myapp
COPY --from=build /src/etc/config.yaml /etc/config/config.yaml
#CMD ["/bin/sh"]

EXPOSE 3000

CMD ["/src/bin/myapp", "-config=/etc/config/config.yaml"]