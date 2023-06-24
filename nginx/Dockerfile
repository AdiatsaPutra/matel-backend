FROM nginx:latest

# Install additional packages
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    liblua5.1-0 \
    lua-nginx-module


# Copy nginx configuration file
# COPY ./nginx.conf /etc/nginx/nginx.conf
