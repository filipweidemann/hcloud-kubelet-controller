# Fully automated Kubelet TLS Approvals on Hetzner Cloud

### Disclaimer
This project is not really anything new, it just handles *one* of many slightly annoying tasks for Kubernetes clusters running **explicitly** on Hetzner Cloud.

This means that the controller handles automatic approval/denial of Kubelet CSRs for your cluster using the Hcloud API to verify the individual CSRs spawned by your Kubelets.

I basically just built this for myself to allow proper TLS bootstrapping for my own clusters without the hassle of approving the resulting CSRs myself each time I grow the cluster, renew certs, ... you get the idea. If you have a use case for this controller and want to use this, feel free to. I'll do my best to keep it updated.

That being said, this is not at all feature complete yet. I plan to expand the checks, add metrics and other stuff.

I'd like to also credit the projects that inspired me to build a Hetzner Cloud integration, namely the great [kubelet-csr-approver by Postfinance](https://github.com/postfinance/kubelet-csr-approver). I had some initial issues with my test setup and looking how they did it really helped me. So check them out if you need a more "general purpose" CSR Approval Controller for your cluster, because this one probably stays exclusive to Hetzner Cloud to implement upstream/oob checks :)

## What problems does this solve?
Whenever you bootstrap a Kubernetes cluster using the default `kubeadm` flow and follow the docs, you'll probably end up with a working cluster.

However, once you try to improve it a little bit, e.g. by adding the [metrics-server](https://github.com/kubernetes-sigs/metrics-server) to your cluster, you'll probably end up with some initial head-scratching because your pod will not start properly.
If you take a closer look at the [Requirements section](https://github.com/kubernetes-sigs/metrics-server/blob/master/README.md#requirements) you'll find that (quote):

"Kubelet certificate needs to be signed by cluster Certificate Authority (or disable certificate validation by passing --kubelet-insecure-tls to Metrics Server)"

So basically, allow insecure (self-signed) Kubelet certificates (don't do that) or handle the TLS bootstrapping correctly. This is where this controller kicks in and handles all the verification checks against the Hetzner Cloud API regarding Kubelet CSRs for you.

## Requirements
Bootstrap your Kubernetes cluster [following the official docs](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-certs/#kubelet-serving-certs) to enable proper "signed serving certificates" to be used by your Kubelets.

Make sure that you also install a proper CNI (e.g. Calico) prior to rolling out the controller or set your taints so it is actually allowed to run on a node.

## Deploy (ressource currently missing)

Once you're done with ensuring the requirements are all met, you can deploy the controller
to your cluster.

##### Clone the repo
`git clone https://github.com/filipweidemann/hcloud-kubelet-controller`

##### Create needed ressources
`kubectl apply -f k8s/req/`

##### Create Deployment
Set the HCLOUD_TOKEN environment variable inside `k8s/deploy.yml`, followed by

`kubectl apply -f k8s/deploy.yml`

Helm Chart coming soon...

