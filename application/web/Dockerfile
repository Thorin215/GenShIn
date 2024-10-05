FROM node:16.17.1-alpine AS web
WORKDIR /root/togettoyou
COPY package*.json ./
RUN npm config set registry https://registry.npmmirror.com/ && npm install
COPY . .
RUN npm run build:prod

FROM nginx:stable-alpine
COPY --from=web /root/togettoyou/dist /usr/share/nginx/html
COPY default.conf /etc/nginx/conf.d/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
