FROM node:16.20-alpine

RUN mkdir -p /var/www/dockerize-nuxt/nuxt-app
WORKDIR /var/www/dockerize-nuxt/nuxt-app

# Install required packages
RUN apk add --update --no-cache git

RUN git clone https://github.com/AdiatsaPutra/matel-frontend . --depth=1

RUN npm install

RUN npm run build

EXPOSE 3000

ENV NUXT_HOST=0.0.0.0

ENV NUXT_PORT=3000

CMD [ "npm", "start" ]
