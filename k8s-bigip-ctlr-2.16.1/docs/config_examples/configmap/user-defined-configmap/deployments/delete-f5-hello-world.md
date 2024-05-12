Below shell script deletes `f5-hello-world` related deployments in default namespace

```
#!/bin/bash

kubectl delete -f f5-hello-world-service.yaml
kubectl delete -f f5-hello-world-deployment.yaml
```
