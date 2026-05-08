# 使用 nginx:alpine 托管前端静态资源
FROM nginx:alpine

# 删除默认页面
RUN rm -rf /usr/share/nginx/html/*

# 复制构建产物
COPY dist/ /usr/share/nginx/html/

# 复制自定义 nginx 配置（支持前端路由）
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
