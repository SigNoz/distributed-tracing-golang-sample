# Base image
FROM nginx:1.23.4-alpine

# Set the working directory
WORKDIR /usr/share/nginx/html

# Copy the static files to the image
COPY frontend .

# Move configuration files to the correct location and cleanup unnecessary files
RUN mv nginx.conf /etc/nginx/conf.d/default.conf && \
    rm -rf Dockerfile

# Expose the port used by Nginx
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]