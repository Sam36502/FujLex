FROM ubuntu:latest

COPY bin/fujlex /app/fujlex
COPY static/ /app/static/
COPY tmpl/ /app/tmpl
#RUN ldconfig

EXPOSE 1919

WORKDIR /app
CMD ["./fujlex"]
