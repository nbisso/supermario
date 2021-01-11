#!/bin/sh

FRONTEND="admin-frontend:0.0.1"
AUTH_SERVICE="auth-service:0.0.1"
USERS_BACKEND="users-backend:0.0.1"
NGINX="mario-nginx:0.0.1"


#minikube start
#eval $(minikube docker-env)
cd ./admin-frontend && docker build -t $FRONTEND  .
cd ..
cd ./auth-service && docker build -t $AUTH_SERVICE  .
cd ..
cd ./users-backend && docker build -t $USERS_BACKEND .
cd ..
cd ./nginx && docker build -t $NGINX .
cd ..

minikube kubectl -- delete -R -f .
minikube kubectl -- apply -R -f .