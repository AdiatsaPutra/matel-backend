FROM nginx:latest

# Install additional packages
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    liblua5.1-0 \
    wget \
    gnupg \
    ca-certificates

RUN wget -O - https://openresty.org/package/pubkey.gpg | sudo apt-key add -

RUN codename=`grep -Po 'VERSION="[0-9]+ \(\K[^)]+' /etc/os-release`

RUN echo "deb http://openresty.org/package/debian $codename openresty" \
    | sudo tee /etc/apt/sources.list.d/openresty.list

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    openresty


# Copy nginx configuration file
# COPY ./nginx.conf /etc/nginx/nginx.conf
