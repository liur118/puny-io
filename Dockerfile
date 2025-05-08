# 第一阶段：前端构建
FROM docker.1ms.run/node:20-alpine AS frontend
WORKDIR /app/ui
COPY ui/package.json ./
RUN npm install --registry=https://mirrors.cloud.tencent.com/npm/
COPY ui .
RUN npm run build

# 第二阶段：后端构建
FROM docker.1ms.run/golang:1.23-alpine AS backend
WORKDIR /app
COPY . .
COPY --from=frontend /app/ui/dist ./ui/dist
RUN apk add --no-cache git --repository=http://mirrors.aliyun.com/alpine/v3.18/main && \
    export GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 go build -o puny-io .

# 第三阶段：极简运行环境
FROM docker.1ms.run/alpine:latest
WORKDIR /app
COPY --from=backend /app/puny-io /app/puny-io
COPY --from=backend /app/conf /app/conf
COPY --from=backend /app/ui/dist /app/ui
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/app/puny-io"] 