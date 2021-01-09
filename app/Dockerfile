FROM python:3.8-slim-buster

EXPOSE 8080
WORKDIR /app

COPY . .
RUN apt-get update && \
    apt-get install --no-install-recommends git --yes && \
    pip install -r requirements.txt

ENTRYPOINT gunicorn main:app --bind 0.0.0.0:8080