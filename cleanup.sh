#!/bin/bash

gcloud compute firewall-rules delete nginx
gcloud compute forwarding-rules delete nginx
gcloud compute instance-groups managed delete nginx
gcloud compute instance-templates delete nginx
gcloud compute target-pools delete nginx
gcloud compute http-health-checks delete nginx
gcloud compute addresses delete nginx

kubectl delete deployments nginx
kubectl delete hpa nginx
kubectl delete svc nginx
