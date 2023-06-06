FROM node:latest

ENV HOST 0.0.0.0

# Install node modules
RUN npm install -g nodemon

# Create app directory
RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app/

COPY . .

ONBUILD RUN npm install

EXPOSE 3000

RUN npm run build
# Install app dependencies
CMD [ "npm", "run", "start"]
