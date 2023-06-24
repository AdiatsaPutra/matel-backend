FROM openresty/openresty:1.21.4.1-0-alpine

# Install LuaRocks
RUN apk --no-cache add curl \
    && curl https://luarocks.org/releases/luarocks-3.7.0.tar.gz | tar xz \
    && cd luarocks-3.7.0 \
    && ./configure \
    && make build \
    && make install \
    && cd .. \
    && rm -rf luarocks-3.7.0

# Install Lua module (resty.http)
RUN luarocks install lua-resty-http