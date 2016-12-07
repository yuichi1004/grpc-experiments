## Add secrets

```bash
$ cd creds
$ kubectl create secret generic creds --from-file=fibocert.pem --from-file=fibokey.pem --from-file=tokencert.pem --from-file=tokenkey.pem
```

## Create RC

```bash
kubectl create -f ./token_rc.yaml
kubectl create -f ./fibo_rc.yaml
```

## Expose

```bash
$ kubectl expose rc token --port=443 --target-port=443 --type=LoadBalancer
$ kubectl expose rc fibo --port=443 --target-port=443 --type=LoadBalancer
```
