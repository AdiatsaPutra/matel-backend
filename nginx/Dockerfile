FROM openresty/openresty:1.21.4.1-0-alpine

# Install LuaRocks
RUN apk --no-cache add curl build-base unzip
RUN curl -LO https://luarocks.org/releases/luarocks-3.9.2.tar.gz
RUN tar xzf luarocks-3.9.2.tar.gz
RUN cd luarocks-3.9.2 && ./configure && make build && make install
RUN rm -rf luarocks-3.9.2 luarocks-3.9.2.tar.gz

# Install Lua module (resty.http)
RUN luarocks install lua-resty-http