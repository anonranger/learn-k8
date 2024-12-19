#### Implementing Pod Security Standards (PSS) with Pod Security Admission (PSA)
*PSA is a built in admission controller that implements the the security controls/policies outlined in the PSS.*

> [!NOTE]
> *Kubernetes Admission Controllers: an Admission Controller is a piece of code that intercepts requests to the Kubernetes API server before they are persisted into etcd and used to change the state of the cluster. Admission controllers can be of type mutating, validating, or both.*


***PSS define three different policies to broadly cover the security spectrum. These policies are cumulative and range from highly-permisive to highly-restrictive.***

#### Policy levels:
 - Privileged: *Unrestricted policy, providing the wildest possible level of permissions. This policy allows for known privilege escalations.*
 - Baseline: *Minimally restrictive policy which prevents known privilege escalations. Allows the default (minimally specified) pod configuration.*
 - Restricted: *Heavily restricted policy, following current pod hardening best practices.*

***PSA admission controller implements the controls, outlined by the PSS policies, via three modes of operation, which are provided in the following list:***
 - enforce: *Policy violations will cause the pod to be rejected.*
 - audit: *Policy violations trigger the addition of an audit annotation to the event recorded in audit log, but are otherwise allowed.*
 - warn: *Policy violations will trigger a user-facing warning, but are otherwise allowed.*

#### PSA labels for Namespaces
You must configure specific PSA modes and PSS profiles at the Kubernetes Namespace level.
With Kubernetes labels, you can choose which of the predefined PSS levels you want to use for pods in a given Namespace.
The labels you select define what action the PSA takes if a potential violation is detected.

> [!IMPORTANT]
> MODE must be one of `enforce`, `audit`, or `warn`.
> LEVEL must be one of `privileged`, `baseline`, or `restricted`.

 ```bash
 #  PSA and PSS Namespace configurations
 #  pod-security.kubernetes.io/<MODE>: <LEVEL>

 kubectl apply -f pss/psa-default.yml
 ```
> [!NOTE]
> PSA enforce mode prevents pods with PSS violations from being applied, but does not stop higher-level controllers, such as Deployments. In fact, the Deployment will be applied successfully without any indication that the pods failed to be applied. While you can use kubectl to inspect the Deployment object, and discover the failed pods message from the PSA, the user experience could be better. To make the user experience better, multiple PSA modes (audit, enforce, warn) should be used.



#### Network Policy
[Network Policy Editor by - ISOVALENT](https://editor.networkpolicy.io/)
OSI Layer 3&4 Network rules using IP Addresses and Ports, It's a builtin kubernetes feature, but it's implementation part of it lies in the CNI which is installed in your cluster. A CNI like Cilium gives you the capabilitie to create custom resources to control layer 7 (HTTP) traffic.

 ```bash
 kubectl get networkpolicies -A
 ```

> [!NOTE]
> By default: all communication (to and from Pods) is allowed in a kubernetes cluster even across namespaces and each pod gets it's own IP Address. This might pose a challenge in a multi-tenent cluster.

***Network policy component is the way to control traffic flow at the IP address and port level.
It sets rules for Pods to communicate, these are strict rules for inter cluster communication.***

Scenario we are solving for:
 - The FRONTEND sends ingress traffic to the WORKER. The WORKER can receive ingress traffic only from the FRONTEND and send ingress traffic to the DATABASE. The DATABASE can receive ingress traffic only from the WORKER and block all egress traffic.

[Use Network Policies with EKS Auto Mode](https://docs.aws.amazon.com/eks/latest/userguide/auto-net-pol.html)

```bash
# Fetch the VPC ID and CIDR block dynamically
export VPC_ID=$(aws eks describe-cluster --name <your-cluster-name> --query "cluster.resourcesVpcConfig.vpcId" --output text)
CIDR=$(aws ec2 describe-vpcs --vpc-ids $VPC_ID --query "Vpcs[0].CidrBlock" --output text)
```

```bash
    # Create or modify your Frontend NetworkPolicy YAML with the dynamically fetched CIDR

cat <<EOF > network-policy/frontend-network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: frontend-network-policy
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: frontend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - ipBlock:
        cidr: $CIDR
    ports:
    - protocol: TCP
      port: 80
    - protocol: TCP
      port: 443
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: backend
    ports:
    - protocol: TCP
      port: 80
EOF
```

```bash
    #Create or modify your Backend NetworkPolicy YAML with the dynamically fetched CIDR

cat <<EOF > network-policy/backend-network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-network-policy
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - ipBlock:
        cidr: $CIDR
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 80
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: mongo
    ports:
    - protocol: TCP
      port: 27017
EOF
```


 ```bash
 # Enable Network Policy Controller

 kubectl apply -f network-policy/.

 kubectl exec <PodName> -- sh -c 'nc -v <PodIP PORT>'
 ```
