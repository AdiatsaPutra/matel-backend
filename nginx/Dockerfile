ARG RESTY_IMAGE_BASE="debian"
ARG RESTY_IMAGE_TAG="bullseye-slim"

FROM ${RESTY_IMAGE_BASE}:${RESTY_IMAGE_TAG}

# RESTY_DEB_FLAVOR build argument is used to select other
# OpenResty Debian package variants.
# For example: "-debug" or "-valgrind"
ARG RESTY_DEB_FLAVOR="-debug"
ARG RESTY_DEB_VERSION="=1.21.4.1-1~bullseye1"
ARG RESTY_APT_REPO="https://openresty.org/package/debian"
ARG RESTY_APT_PGP="https://openresty.org/package/pubkey.gpg"

LABEL resty_image_base="${RESTY_IMAGE_BASE}"
LABEL resty_image_tag="${RESTY_IMAGE_TAG}"
LABEL resty_apt_repo="${RESTY_APT_REPO}"
LABEL resty_apt_pgp="${RESTY_APT_PGP}"
LABEL resty_deb_flavor="${RESTY_DEB_FLAVOR}"
LABEL resty_deb_version="${RESTY_DEB_VERSION}"


RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ca-certificates \
    gettext-base \
    gnupg2 \
    lsb-base \
    lsb-release \
    software-properties-common \
    wget \
    && wget -qO /tmp/pubkey.gpg ${RESTY_APT_PGP} \
    && DEBIAN_FRONTEND=noninteractive apt-key add /tmp/pubkey.gpg \
    && rm /tmp/pubkey.gpg \
    && DEBIAN_FRONTEND=noninteractive add-apt-repository -y "deb ${RESTY_APT_REPO} $(lsb_release -sc) openresty" \
    && DEBIAN_FRONTEND=noninteractive apt-get remove -y --purge \
    gnupg2 \
    lsb-release \
    software-properties-common \
    wget \
    && DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    openresty${RESTY_DEB_FLAVOR}${RESTY_DEB_VERSION} \
    && DEBIAN_FRONTEND=noninteractive apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /var/run/openresty \
    && ln -sf /dev/stdout /usr/local/openresty${RESTY_DEB_FLAVOR}/nginx/logs/access.log \
    && ln -sf /dev/stderr /usr/local/openresty${RESTY_DEB_FLAVOR}/nginx/logs/error.log

# Install LuaRocks
RUN /usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/bin/luarocks install lua-resty-http
RUN /usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/bin/luarocks install lua-cjson

# Add additional binaries into PATH for convenience
ENV PATH=$PATH:/usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/bin:/usr/local/openresty${RESTY_DEB_FLAVOR}/nginx/sbin:/usr/local/openresty${RESTY_DEB_FLAVOR}/bin
ENV LUA_PATH="/usr/local/openresty${RESTY_DEB_FLAVOR}/site/lualib/?.ljbc;/usr/local/openresty${RESTY_DEB_FLAVOR}/site/lualib/?/init.ljbc;/usr/local/openresty${RESTY_DEB_FLAVOR}/lualib/?.ljbc;/usr/local/openresty${RESTY_DEB_FLAVOR}/lualib/?/init.ljbc;/usr/local/openresty${RESTY_DEB_FLAVOR}/site/lualib/?.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/site/?.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/site/lualib/?/init.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/site/?/init.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/lualib/?.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/lualib/?/init.lua;./?.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/share/luajit-2.1.0-beta3/?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/share/lua/5.1/?.lua;/usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/share/lua/5.1/?/init.lua"
ENV LUA_CPATH="/usr/local/openresty${RESTY_DEB_FLAVOR}/site/lualib/?.so;/usr/local/openresty${RESTY_DEB_FLAVOR}/lualib/?.so;./?.so;/usr/local/lib/lua/5.1/?.so;/usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/lib/lua/5.1/?.so;/usr/local/lib/lua/5.1/loadall.so;/usr/local/openresty${RESTY_DEB_FLAVOR}/luajit/lib/lua/5.1/?.so"

CMD ["/usr/local/openresty${RESTY_DEB_FLAVOR}/bin/openresty", "-g", "daemon off;"]