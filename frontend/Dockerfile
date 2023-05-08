FROM artifactory.cloud.cms.gov/docker/nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY src/* /var/www/kl/

ENV SERVER_NAME example.com


RUN apk add --no-cache libc6-compat && \
  sed -i 's/example.com/'${SERVER_NAME}'/' /etc/nginx/nginx.conf

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
