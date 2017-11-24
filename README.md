# kube-default-http-backend

HTTP backend to be used as default for Ingress controller. Writen to match [requirements](https://github.com/kubernetes/ingress-nginx#requirements) imposed by [nginx ingress controller](https://github.com/kubernetes/ingress-nginx). Based on [this sample code](
https://github.com/kubernetes/ingress-nginx/tree/master/images/custom-error-pages).

Tested only with nginx ingress controller.

## Usage

1. Install [deployment](deployment.yaml)
2. Expose it as a service (`kubectl expose default-http-backend`)
3. Point your nginx ingress controller to the service (set `--default-backend-service=NAMESPACE/default-http-backend`)
4. [Configure](https://github.com/kubernetes/ingress-nginx/blob/master/docs/user-guide/configmap.md#custom-http-errors) which errors controller should pass to default backend
