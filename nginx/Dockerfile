FROM nginx:1.18.0-alpine
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d
COPY secured-location.conf /etc/nginx
EXPOSE 9999
CMD ["nginx", "-g", "daemon off;"]