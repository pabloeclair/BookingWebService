FROM node:20

WORKDIR /app

COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

CMD ["sh", "-c", "sleep 13 && npm run serve"]
