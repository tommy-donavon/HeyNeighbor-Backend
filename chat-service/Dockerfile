FROM node:16.13.1-alpine3.14 as build

WORKDIR /usr

COPY package.json ./

COPY tsconfig.json ./

COPY src ./src

RUN ls -a

RUN npm install

RUN npm run build

## this is stage two , where the app actually runs

FROM node:16.13.1-alpine3.14

WORKDIR /usr

COPY package.json ./

RUN npm install --only=production

COPY --from=build /usr/dist .

# RUN npm install pm2 -g


CMD ["node", "--experimental-modules","app.js"]