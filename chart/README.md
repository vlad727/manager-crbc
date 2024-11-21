## Manager for kubernetes namespaces
Application provide web UI and allow creates cluster role bindings for service account also in annotate cluster role binding like  "requester: <logged-username>"
### How to build application
```
docker build -t manager-crbc.
```
### How to deploy it to k8s

### Install with Helm
```
# Install helm before
brew install helm 
# Clone this project and install by helm
helm repo add postinstall https://registry.apps.k8s.ose-prod.solution.sbt/chartrepo/postinstall --username LOGIN --password TOKEN
helm pull postinstall/manager-crbc --untar
helm install crbc chart/ 

# Check installation result
kubectl get pods -n manager-crbc
open ingress https://manager-crbc.apps.<your-cluster--fqdn>>.solution.sbt/
create cluster role binding via web
```
