FROM node:latest

COPY . .

# Install node modules
RUN npm install

# Create app directory
RUN mkdir -p /app
WORKDIR /app

EXPOSE 3000

RUN npm run build
# Install app dependencies

CMD [ "npm", "run", "start" ]
