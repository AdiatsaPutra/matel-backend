FROM node:latest

ENV HOST 0.0.0.0

COPY . .

# Install node modules
RUN npm install

# Create app directory
RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app/

EXPOSE 8080

RUN npm run build
# Install app dependencies

CMD [ "npm", "run", "start" ]
