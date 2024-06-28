# Local Testing Environment for OAI with O1-Telnet
Follow the [Installation Guide](https://docs.nephio.org/docs/guides/install-guides/) to install Nephio

Now, we need to follow the [following exercise](https://docs.nephio.org/docs/guides/user-guides/exercise-2-oai/) which we term as Master Excerises.

### Step-1: Setting up infrastructure
Follow the steps 1 and 2 of [Master Excercise](https://docs.nephio.org/docs/guides/user-guides/exercise-2-oai/)

### Step-2: Creating a Local operator-image and pushing it to Edge and Regional Cluster:
``` bash
# Clone the Repo
cd oai/
docker build . -t local-ran-operator:v0.1
kind load docker-image local-ran-operator:v0.1  -n edge
kind load docker-image local-ran-operator:v0.1  -n regional
``` 

### Step-3. Setting up Local Blueprint Repo:
 1. Create a new Repo at the Gitea-Cluster (172.18.0.200:3000) (Repo-name: 'oai-ran-bp' is recommended)
 2. Register the Repo to KPT
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
  The above script will create the following blueprint packages with updated values: 'oai-ran-operator', 'pkg-example-cucp-bp', 'pkg-example-cuup-bp', 'pkg-example-du-bp'.
  Make sure to check your repo, before moving forward.

### Step-4: Making sure your PackageVariant point to your local Bp-Repo
```bash
 cd test-infra/local-bp/packageVariants/
 cp -r . $HOME/test-infra/e2e/tests/oai/
```
Note: If your repo-name is not 'oai-ran-bp', then update the upstream-repo of packageVariant accordingly.

### Step-5: Deploy Core and Ran
Follow the steps 3-6 of [Master Excercise](https://docs.nephio.org/docs/guides/user-guides/exercise-2-oai/)

### Step-6: Validate the deployment by ascessing o1-telnet
```bash
 sudo apt update && sudo apt install netcat
 O1_IP=$(kubectl get svc oai-gnb-du-o1-telnet-lb -n oai-ran-du --context edge-admin@edge -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
 echo o1 stats | nc -N  $O1_IP 9090
```

### Step-7: Deploy UE (20 MHz):
The current ue-image is not able to connect to the du using service (oai-gnb-du.oai-ran-du), Therefore, we provide the du-ip to ue manually as workaround: (Todo: Ideal way to resolve)
```bash
DU_POD=$(kubectl get pods -n oai-ran-du --context edge-admin@edge  -l app.kubernetes.io/name=oai-gnb-du  -o jsonpath='{.items[*].metadata.name}')
DU_IP=$(kubectl get pods $DU_POD -n oai-ran-du --context edge-admin@edge -o jsonpath="{.status.podIP}")
echo $DU_IP
```
Copy the du-ip to the "uedeployment.yaml" under the Additional-Option "--rfsimulator.serveraddr" to both 20Mhz and 40Mhz UE.

#### Deploy the ue
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
### Step-8: Bandwidth Reconfigure Procedure (20 Mhz to 40 Mhz) using o1-telnet
```bash
O1_IP=$(kubectl get svc oai-gnb-du-o1-telnet-lb -n oai-ran-du --context edge-admin@edge -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo o1 stop_modem | nc -N $O1_IP 9090
echo o1 bwconfig 40 | nc -N $O1_IP 9090
echo o1 start_modem | nc -N $O1_IP 9090
echo o1 stats | nc -N $O1_IP 9090 
```
After Reconfiguration, Connect the 40 Mhz ue using Step-7 and validate the Bandwidth-reconfiguration using Ping-Test
