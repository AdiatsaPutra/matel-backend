FROM nginx:latest

# Install dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    build-essential \
    libpcre3-dev \
    libssl-dev \
    libxml2-dev \
    libxslt1-dev \
    zlib1g-dev \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Download and extract OpenResty
RUN /bin/bash -c "wget https://openresty.org/download/openresty-1.19.3.2.tar.gz \
    && tar -xzf openresty-1.19.3.2.tar.gz \
    && rm openresty-1.19.3.2.tar.gz"

# Build and install OpenResty
RUN cd openresty-1.19.3.2 \
    && ./configure --with-http_stub_status_module --with-http_realip_module --add-module=/path/to/ngx_http_lua_module \
    && make \
    && make install

# Copy nginx configuration file
# COPY ./nginx.conf /etc/nginx/nginx.conf
