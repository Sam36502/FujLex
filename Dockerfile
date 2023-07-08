FROM ubuntu:latest

ADD fujlex /fujlex
ADD static /static
ADD tmpl /tmpl
RUN ldconfig

EXPOSE 1919
CMD ["/fujlex"]
