FROM node:latest

# Install node modules
RUN npm install -g nodemon

# Create app directory
RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app/

ONBUILD RUN npm install

EXPOSE 3000

ENV HOST 0.0.0.0
# Install app dependencies
CMD [ "npm", "run", "buildStart"]
