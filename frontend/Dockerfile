# Base image
FROM nginx:1.23.4-alpine

# Copy the static files to the image
COPY frontend /usr/share/nginx/html

# Move configuration files to the correct location and cleanup unnecessary files
RUN mv frontend/nginx.conf /etc/nginx/conf.d/default.conf && \
    mv frontend/src/docker-config.js frontend/src/config.js && \
    rm -rf /usr/share/nginx/Dockerfile

# Expose the port used by Nginx
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]