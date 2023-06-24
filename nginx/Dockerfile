FROM openresty/openresty:1.21.4.1-0-focal

# Install additional packages
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    libssl-dev \
    perl \
    unzip

# Install LuaRocks
RUN /usr/local/openresty/luajit/bin/luarocks install lua-resty-http
