FROM ubuntu:14.04
ADD ./release
RUN pip install -r requirements.txt
CMD go_shortlink release/redis-3.0.5/src/redis-cli
