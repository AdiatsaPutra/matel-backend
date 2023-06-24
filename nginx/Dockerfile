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
    wget \
    liblua5.1-dev \
    && rm -rf /var/lib/apt/lists/*

# Download and extract OpenResty
RUN wget https://openresty.org/download/openresty-1.19.3.2.tar.gz \
    && tar -xzf openresty-1.19.3.2.tar.gz \
    && rm openresty-1.19.3.2.tar.gz

# Download and move ngx_http_lua_module
RUN git clone https://github.com/openresty/lua-nginx-module.git \
    && mv lua-nginx-module openresty-1.19.3.2/

# Build and install OpenResty with ngx_http_lua_module
RUN cd openresty-1.19.3.2 \
    && ./configure --with-http_stub_status_module --with-http_realip_module --add-module=./lua-nginx-module \
    && make \
    && make install


# Copy nginx configuration file
# COPY ./nginx.conf /etc/nginx/nginx.conf
