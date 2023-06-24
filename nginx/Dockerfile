FROM openresty/openresty:1.21.7.1-alpine

# Install LuaRocks
RUN apk --no-cache add curl
RUN curl -LO https://luarocks.org/releases/luarocks-3.7.0.tar.gz
RUN tar xzf luarocks-3.7.0.tar.gz
RUN cd luarocks-3.7.0 && ./configure && make build && make install
RUN rm -rf luarocks-3.7.0 luarocks-3.7.0.tar.gz

# ...



# Install Lua module (resty.http)
RUN luarocks install lua-resty-http