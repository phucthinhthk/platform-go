FROM node:23-alpine

RUN apk update
# openapi-generator-cli
RUN apk add --no-cache openjdk21-jre-headless

RUN npm install -g npm
# openapi cli
RUN npm install -g @openapitools/openapi-generator-cli
# 
RUN openapi-generator-cli version

# Thêm openapi2aspida vì bạn muốn chạy make codegen
RUN npm install -g openapi2aspida typescript ts-node

WORKDIR /app

# Giữ container chạy ngầm
CMD ["sh"]
