#!/bin/bash
make
docker build -t registry-vpc.cn-hangzhou.aliyuncs.com/leo123/leo-controller:v1 .
docker push registry-vpc.cn-hangzhou.aliyuncs.com/leo123/leo-controller:v1
kubectl delete -f deploy/leo-controller-deployment.yml
kubectl apply -f deploy/leo-controller-deployment.yml
docker images | grep none | awk '{print $3}' | xargs docker rmi
# make build
# make push
# make apply