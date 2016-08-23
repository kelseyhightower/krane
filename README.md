# Krane

Convert Google Compute Engine (GCE) autoscaling instance groups to Kubernetes autoscaling deployments.

## Status

Proof of concept.

## Install

Download the latest release from the [releases page](https://github.com/kelseyhightower/krane/releases).

## Usage

### Create a Kubernetes cluster

```
gcloud container clusters create cluster0
```

### Create an Autoscaling GCE Instance Group

Instances will leverage the Google Container VM and the Kubernetes Kubelet to start applications from Pod manifests. See the cloud-config.yaml for more details.

```
gcloud compute addresses create nginx
```

```
gcloud compute http-health-checks create nginx
```

```
gcloud compute target-pools create nginx --health-check nginx
```

```
gcloud compute instance-templates create nginx \
  --image-family gci-stable \
  --image-project google-containers \
  --machine-type n1-standard-1 \
  --metadata-from-file user-data=cloud-config.yaml
```

```
gcloud compute instance-groups managed create nginx \
  --base-instance-name nginx \
  --target-pool nginx \
  --size 2 \
  --template nginx
```

```
gcloud compute instance-groups managed set-autoscaling nginx \
  --max-num-replicas 10 \
  --target-cpu-utilization 0.25 \
  --cool-down-period 90
```

```
gcloud compute forwarding-rules create nginx \
  --ports 80 \
  --address $(gcloud compute addresses describe nginx --format 'value(address)') \
  --target-pool nginx
```

```
gcloud compute firewall-rules create nginx --allow tcp:80
```

### Convert a GCE Instance Group to Kubernetes

```
krane -instance-group nginx | kubectl create -f -
```

```
deployment "nginx" created
horizontalpodautoscaler "nginx" created
service "nginx" created
```

> To view the Kubernetes configs omit the pipe to kubectl: `krane -instance-group nginx`

#### Results

The Pod manifest from the instance group cloud-init metadata is converted to a Kubernetes deployment.

```
kubectl get deployments
```
```
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx     2         2         2            2           3m
```

```
kubectl get pods
```
```
NAME                     READY     STATUS    RESTARTS   AGE
nginx-2115999728-c5s11   1/1       Running   0          3m
nginx-2115999728-kbfj0   1/1       Running   0          3m
```

The autoscaling policy from the instance group is converted to a Kubernetes horizontal pod autoscaler.

```
kubectl get hpa
```
```
NAME      REFERENCE          TARGET    CURRENT   MINPODS   MAXPODS   AGE
nginx     Deployment/nginx   25%       0%        2         10        3m
```

If the instance group is part of a target pool which is referenced by a network load balancer as Kubernetes service will be created.

```
kubectl get svc
```
```
NAME         CLUSTER-IP     EXTERNAL-IP    PORT(S)   AGE
nginx        XX.XX.XXX.XX   XXX.XXX.XX.X   80/TCP    3m
```
