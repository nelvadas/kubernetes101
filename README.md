# From Monolith to Cloud Native Applications with Docker and Kubernetes  101

The purpose of this lab, is to quickly show how to  migrate a monolith back end service to Kubernetes.
The lab is organized in 3 parts
* the first part show the legacy hellok8s service in Go
* the second part shows how to containerize the service using docker
* the last part focuses on running the service on a minukube cluster.



## STEP 1:  Monolith helloword service

Navigate in source code hello/main.go

The web service just print a welcome message in the server logs (DEBUG mode) .
The same message is returned to the caller in INFO

![main.go ](https://github.com/nelvadas/kubernetes101/blob/master/step3/maingoclass.png "helloservice ")



```
step1$ go build hello/main.go
step1$ go run hello/main.go
Server starting on :8080..
```


Test the monolith application from your localhost using curl/httpie

```
$ http localhost:8080
HTTP/1.1 200 OK
Content-Length: 56
Content-Type: text/plain; charset=utf-8
Date: Sun, 22 Jul 2018 18:28:44 GMT

[2018.07.22 20:28:44] INFO Docker and Kubernetes 101   


```

the server log print the customer result in DEBUG MODE
```
Server starting on :8080..
[2018.07.22 20:28:33] DEBUG Docker and Kubernetes 101
[2018.07.22 20:28:44] DEBUG Docker and Kubernetes 101
```

How run the following microservice in a Kubernetes environment?
The containerization needs two steps: Dockerize and Kubifying:)



## STEP 2: Containerizing with docker

Build the docker image

```
$ cd step2
$ docker build --no-cache -t hellok8s:1.0.0 .
Sending build context to Docker daemon  3.584kB
Step 1/5 : FROM golang:latest
 ---> 4e611157870f
Step 2/5 : ADD .  /go/src/golang/
 ---> b4710a90d4cc
Step 3/5 : RUN go install golang/hello
 ---> Running in dcef1971a64f
Removing intermediate container dcef1971a64f
 ---> dd26b5026650
Step 4/5 : ENTRYPOINT /go/bin/hello
 ---> Running in 4cbfea221fc4
Removing intermediate container 4cbfea221fc4
 ---> 6a136fc4afaf
Step 5/5 : EXPOSE 8080
 ---> Running in ce0f197a5a74
Removing intermediate container ce0f197a5a74
 ---> c39d1169e137
Successfully built c39d1169e137
Successfully tagged hellok8s:1.0.0
 ```

*Q : How to add symbolic links in a docker images*

check the generated image

```
REPOSITORY                                                         TAG                 IMAGE ID            CREATED             SIZE
hellok8s                                                           1.0.0               c39d1169e137        2 minutes ago       801MB
...
golang                                                             latest              4e611157870f        3 weeks ago         794MB
```

Run the docker container

```
$ docker run --rm --name hellocontainer1  -d -p 7070:8080 hellok8s:1.0.0
d92844fd071359567088b2d7d5af12f837056e7ebf8a3aeb01506ba3ab2dbc5e
```

```
$ docker ps
CONTAINER ID   IMAGE      COMMAND   CREATED             STATUS             PORTS                    NAMES
d92844fd0713   hellok8s:1.0.0"/bin/sh -c /go/bin/…"   6 sec ago Up 4 sec   0.0.0.0:7070->8080/tcp   hellocontainer1
```

Test the hello microservice from the container endpoint
```
$ http localhost:7070
HTTP/1.1 200 OK
Content-Length: 56
Content-Type: text/plain; charset=utf-8
Date: Sun, 22 Jul 2018 18:47:21 GMT

[2018.07.22 18:47:21] INFO Docker and Kubernetes 101   

```

Display the logs of your container
```
$ docker logs d92844fd0713
Server starting on :8080..
[2018.07.22 18:47:21] DEBUG Docker and Kubernetes 101
```

*use  docker logs -f CONTAINER ID  to see the logs in real time*

Create a second container instance for the hello service

```

$ docker run --rm --name hellocontainer2  -d -p 7071:8080 hellok8s:1.0.0
40fde06d92a8651547132baa2788c5042c0347940dd523a31f17c37f669bb459
```
Check the instance 2
```
$ http localhost:7071
HTTP/1.1 200 OK
Content-Length: 56
Content-Type: text/plain; charset=utf-8
Date: Sun, 22 Jul 2018 18:52:21 GMT

[2018.07.22 18:52:21] INFO Docker and Kubernetes 101
```

Now we can push the image on docker hub to share with others :)

```
echo MYSECRETPASSWORD | docker login -u nelvadas --password-stdin
Login Succeeded
```

Tag and push the image on dockerhub
```

 $ docker tag  hellok8s:1.0.0 nelvadas/hellok8s:1.0.0
$ docker push nelvadas/hellok8s:1.0.0
The push refers to repository [docker.io/nelvadas/hellok8s]
8c515d306fe8: Pushed
a74fc7fd8729: Pushed
07912702a91b: Pushed
89a6ac27aba4: Pushing [=================>                                 ]  123.8MB/358.9MB
ae717a8370e0: Pushing [=============================================>     ]  146.8MB/161.7MB
c30dae2762bd: Pushed
43701cc70351: Pushed
e14378b596fb: Pushed
a2e66f6c6f5f: Pushing [====================================>              ]  73.92MB/100.6MB

1.0.0: digest: sha256:15123d4bd47305979f90f845af9889a59302b22f81739be5611f277da23f3af5 size: 2214
```






How do you enable clustering features, loadbalancing ,
maintain 2 running instances ?
automatically create new instances when the load increase?

We need an orchestrator: KS8



## Orchestrating with Kubernetes

Install a [minikube](https://github.com/kubernetes/minikube) cluster

Start minikube
```
$ minikube start --vm-driver=virtualbox
Starting local Kubernetes v1.10.0 cluster...
Starting VM...
Getting VM IP address...
Moving files into cluster...
...
```

Check the cluster STATUS
```
$ minikube status
minikube: Running
cluster: Running
kubectl: Correctly Configured: pointing to minikube-vm at 192.168.99.100
MBP-de-elvadas:~ enonowog$
```


create a namespace

```
$ kubectl config set-context --cluster=minikube
$ kubectl create namespace hello
```


Run the image in the kubernete cluster

```
$ kubectl run hellok8s --image=nelvadas/hellok8s:1.0.0 --port=8080 -n hello
```

This create a deployment object
```
$ oc get deployment -n hello
NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
hellok8s   1         1         1            1           13s
```

By default each deployment generate 1 pod
```
$ oc get pods -n hello
NAME                        READY     STATUS    RESTARTS   AGE
hellok8s-6d99bd7994-xgzgl   1/1       Running   0          55s
```

To access the service externally , we need to expose it

```
kubectl expose deployment hellok8s --type=LoadBalancer -n hello
service/hellok8s exposed
```
Access the service using minikube command
```
$ minikube service hellok8s -n hello
Opening kubernetes service hello/hellok8s in default browser...
```


Minikube open the service in the web browser, you can reach the service also by curling the provided
service cluster address.


```
$ http http://192.168.99.100:31421
HTTP/1.1 200 OK
Content-Length: 56
Content-Type: text/plain; charset=utf-8
Date: Sun, 22 Jul 2018 19:30:06 GMT

[2018.07.22 19:30:06] INFO Docker and Kubernetes 101

```

You can see in the pod logs that the request was served
```
$ kubectl get pods -n hello
NAME                        READY     STATUS    RESTARTS   AGE
hellok8s-6d99bd7994-xgzgl   1/1       Running   0          17m

MBP-de-elvadas:~ enonowog$ kubectl logs hellok8s-6d99bd7994-xgzgl -n hello
Server starting on :8080..
[2018.07.22 19:27:42] DEBUG Docker and Kubernetes 101
[2018.07.22 19:27:42] DEBUG Docker and Kubernetes 101
[2018.07.22 19:28:44] DEBUG Docker and Kubernetes 101
[2018.07.22 19:30:06] DEBUG Docker and Kubernetes 101
```

Now we can scale the service and still use the provided url to access both instances


```
$ kubectl scale deployment/hellok8s --replicas=2 -n hello
deployment.extensions/hellok8s scaled
MBP-de-elvadas:~ enonowog$ oc get pods -n hello
NAME                        READY     STATUS    RESTARTS   AGE
hellok8s-6d99bd7994-d4jcg   1/1       Running   0          10s
hellok8s-6d99bd7994-xgzgl   1/1       Running   0          20m

```

A new pod is created and the trafic should now be loadbalanced between the two pods


Open logs of each pod in a separate console tabs using

```
$ kubectl logs -f hellok8s-6d99bd7994-d4jcg -n hello
$ kubectl logs -f hellok8s-6d99bd7994-xgzgl -n hello
```

Send 5 requests


```
$ ab -n 5 -v 1 http://192.168.99.100:31421/
```

![Loadbalancing ](https://github.com/nelvadas/kubernetes101/blob/master/step3/loadbalancing.png "kubectl logs -f ")


we have so far used the kubectl command, it is time to see how the minikube console looks like.

```
$ minikube dashboad
```

![Dashboard ](https://github.com/nelvadas/kubernetes101/blob/master/step3/minikubedashboard.png "minikube console")
