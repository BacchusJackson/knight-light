FROM artifactory.cloud.cms.gov/docker/nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY src/* /var/www/kl/

ENV SERVER_NAME example.com
ARG COUNTER_API
ARG CHAT_API


RUN apk add --no-cach libc6-compat && \
  sed -i 's/example.com/'${SERVER_NAME}'/' /etc/nginx/nginx.conf && \
  echo "const COUNTER_API_ENDPOINT = ${COUNTER_API};" > /var/www/kl/config.js && \
  echo "const CHAT_API_ENDPOINT = ${CHAT_API};" >> /var/www/kl/config.js
    

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
