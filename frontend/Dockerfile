FROM node:20

WORKDIR /app

COPY package*.json ./
RUN npm config set registry https://registry.npmmirror.com/
RUN npm ci
COPY . .
RUN npm run build

CMD ["sh", "-c", "sleep 13 && npm run serve"]
