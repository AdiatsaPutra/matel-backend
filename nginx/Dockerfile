FROM openresty/openresty:1.21.4.1-0-alpine

# Install LuaRocks
RUN apk --no-cache add curl build-base unzip git
RUN curl -LO https://luarocks.org/releases/luarocks-3.7.0.tar.gz
RUN tar xzf luarocks-3.7.0.tar.gz
RUN cd luarocks-3.7.0 && ./configure && make build && make install
RUN rm -rf luarocks-3.7.0 luarocks-3.7.0.tar.gz

# Install Lua module (resty.http)
ENV ROCKS_SERVER="https://luarocks.org/mirror"
RUN luarocks install lua-resty-http2
