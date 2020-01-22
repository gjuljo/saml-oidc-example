# Use the standard Nginx image from Docker Hub
FROM nginx
RUN apt-get update
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 8443
CMD ["nginx", "-g", "daemon off;"]