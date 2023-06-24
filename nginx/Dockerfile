# Base image menggunakan image Nginx resmi
FROM nginx:latest

# Update paket dan instal dependensi yang diperlukan
RUN apt-get update && apt-get install -y libpcre3-dev libssl-dev zlib1g-dev wget

# Unduh sumber modul Lua
RUN wget https://github.com/openresty/lua-nginx-module/archive/v0.10.20.tar.gz && tar -xzvf v0.10.20.tar.gz

# Unduh sumber LuaJIT
RUN wget http://luajit.org/download/LuaJIT-2.1.0-beta3.tar.gz && tar -xzvf LuaJIT-2.1.0-beta3.tar.gz

# Kompilasi dan instalasi LuaJIT
RUN cd LuaJIT-2.1.0-beta3 && make && make install

# Kompilasi dan instalasi modul Lua pada Nginx
RUN cd lua-nginx-module-0.10.20 && cp -r * /etc/nginx/modules/

# Tambahkan modul Lua ke dalam konfigurasi Nginx
RUN echo "load_module modules/ngx_http_lua_module.so;" >> /etc/nginx/nginx.conf
