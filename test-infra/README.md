# Local Testing Environment

Follow the [Installation Guide](https://docs.nephio.org/docs/guides/install-guides/) to install Nephio

Now, follow the [OAI exercise](https://docs.nephio.org/docs/guides/user-guides/usecase-user-guides/exercise-2-oai/)

### Step-1: Setup the infrastructure

Follow the steps 1 and 2 of [OAI exercise](https://docs.nephio.org/docs/guides/user-guides/usecase-user-guides/exercise-2-oai/)

### Step-2: Create a Local operator-image and push it to Edge and Regional Cluster

``` bash
# Clone the Repo
cd oai/
docker build . -t local-ran-operator:v0.1
kind load docker-image local-ran-operator:v0.1  -n edge
kind load docker-image local-ran-operator:v0.1  -n regional
``` 

### Step-3: Setup Local Blueprint Repository

1. Create a new repository in the Gitea-Cluster (172.18.0.200:3000) (Repo-name: `oai-ran-bp`)

2. Register the repository via KPT

   ``` bash
   kpt alpha repo register \
  --namespace default \
  --repo-basic-username=<username_for_gitea> \
  --repo-basic-password=<passowrd_for_gitea> \
  <repo-url>
   ```
3. Create and push blueprints to local as per the [PR](https://github.com/nephio-project/catalog/pull/41)

   ```bash
   cd test-infra/local-bp/
   ./create-bp-package.sh
   ```
  The above script will create the following blueprint packages with updated values: `oai-ran-operator`, `pkg-example-cucp-bp`, `pkg-example-cuup-bp`, `pkg-example-du-bp`. Make sure to check your repository, before moving forward.

### Step-4: Make sure this PackageVariant points to the local Blueprint repository

```bash
 cd test-infra/local-bp/packageVariants/
 cp -r . $HOME/test-infra/e2e/tests/oai/
```
**Note**: If the repo-name is not `oai-ran-bp`, then update the upstream-repo of packageVariant accordingly.

### Step-5: Deploy Core and RAN

Follow the steps 3-6 of [OAI Exercise](https://docs.nephio.org/docs/guides/user-guides/usecase-user-guides/exercise-2-oai/)

### Step-6: Validate the deployment by accessing telnet

```bash
 sudo apt update && sudo apt install netcat
 TELNET_IP=$(kubectl get svc oai-gnb-du-telnet-lb -n oai-ran-du --context edge-admin@edge -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
 echo o1 stats | nc -N  $TELNET_IP 9090
```

### Step-7: Deploy UE (20 MHz):

```bash
cd test-infra/oai-ue/
kubectl apply -f namespace.yaml --context edge-admin@edge
kubectl apply -f 20Mhz/. --context edge-admin@edge
```

#### Run the Ping Test

```bash
UE_POD=$(kubectl get pods -n oai-ue --context edge-admin@edge  -l app.kubernetes.io/name=oai-nr-ue -o jsonpath='{.items[*].metadata.name}')
UPF_POD=$(kubectl get pods -n oai-core --context=edge-admin@edge -l workload.nephio.org/oai=upf -o jsonpath='{.items[*].metadata.name}')
UPF_tun0_IP_ADDR=$(kubectl exec -it $UPF_POD -n oai-core -c upf-edge --context edge-admin@edge -- ip -f inet addr show tun0 | sed -En -e 's/.*inet ([0-9.]+).*/\1/p')
kubectl exec -it $UE_POD -n oai-ue --context edge-admin@edge -- ping -c 3 $UPF_tun0_IP_ADDR
```
### Step-8: Bandwidth Reconfigure Procedure (20 Mhz to 40 Mhz) using telnet

```bash
TELNET_IP=$(kubectl get svc oai-gnb-du-telnet-lb -n oai-ran-du --context edge-admin@edge -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo o1 stop_modem | nc -N $TELNET_IP 9090
echo o1 bwconfig 40 | nc -N $TELNET_IP 9090
echo o1 start_modem | nc -N $TELNET_IP 9090
echo o1 stats | nc -N $TELNET_IP 9090 
```

After Reconfiguration, Connect the 40 Mhz ue using Step-7 and validate the Bandwidth-reconfiguration using Ping-Test
