# sudo docker build . -t="micro:latest"
# sudo docker run -d -p 9001:9001 -p 7000:7000  micro:latest
FROM ubuntu-env:latest

RUN mkdir -p /app
ADD ./micro /app
COPY micro.conf /etc/supervisor/conf.d

WORKDIR /app

EXPOSE 7000
EXPOSE 9001

CMD ["supervisord",  "-c", "/etc/supervisor/supervisord.conf"]
