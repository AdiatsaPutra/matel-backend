FROM openresty/openresty:1.21.4.1-0-focal

# Install LuaRocks
RUN /usr/local/openresty/luajit/bin/luarocks install lua-resty-http

CMD ["/usr/local/openresty/bin/openresty", "-g", "daemon off;"]