# Helm CRD Rollback Test

This guide provides step-by-step instructions for testing the installation, upgrade, and rollback of a leader election CRD and operator. The leader election example is based on [java-operator-sdk](https://github.com/operator-framework/java-operator-sdk). The process involves installing two Helm packages: first the CRD, then the leader election operator.

## Prerequisites

- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed
- [Helm](https://helm.sh/docs/intro/install/) installed

## Setup

1. Create a Kubernetes cluster using kind:

```bash
kind create cluster
```

## Installation Steps

### Install v1alpha1 Version

1. Install the CRD v1alpha1 version:

   ```bash
   helm install leader-election-crd crd-v1alpha1/chart
   ```

2. Install the leader election application v1alpha1:

   ```bash
   helm install leader-election leader-election-v1alpha1/chart
   ```

### Upgrade to v1alpha2 Version

1. Install the CRD conversion webhook:

   ```bash
   helm install leader-election-crd-conversion crd-conversion/chart
   ```

   > **Note:** This conversion webhook is essential when the CRD's `conversion.strategy` is set to `Webhook`, to support the conversion between v1alpha1 and v1alpha2 CRs.
   > Optionally, you can skip this step by setting the conversion strategy to `None`, and uncomment this line in the CRD for `v1alpha1`:

   ```yaml
   x-kubernetes-preserve-unknown-fields: true
   ```

   > If this step is skipped or if the conversion strategy is set to `None` without turning on `x-kubernetes-preserve-unknown-fields` for `v1alpha1`, you will encounter the following error after upgrading in step 3:
   >
   > ```console
   > [ERROR] Error during event processing ExecutionScope{ resource id: ResourceID{name='leader1', namespace='default'}, version: 80343}
   > io.fabric8.kubernetes.client.KubernetesClientException: Failure executing: PATCH at: https://10.96.0.1:443/apis/sample.operator.javaoperatorsdk.io/v1alpha2/namespaces/default/leaderelections/leader1/status?fieldManager=leaderelectiontestreconciler&force=true. Message: failed to prune fields: failed to convert merged object to last applied version: .status.addedField: field not declared in schema. Received status: Status(apiVersion=v1, code=500, details=null, kind=Status, message=failed to prune fields: failed to convert merged object to last applied version: .status.addedField: field not declared in schema, metadata=ListMeta(_continue=null, remainingItemCount=null, resourceVersion=null, selfLink=null, additionalProperties={}), reason=null, status=Failure, additionalProperties={}).
   > ```

2. Upgrade the CRD to v1alpha2:

   ```bash
   helm upgrade leader-election-crd crd-v1alpha2/chart
   ```

3. Upgrade the leader election application to v1alpha2:

   ```bash
   helm upgrade leader-election leader-election-v1alpha2/chart
   ```

## Rollback Procedure

### Application Rollback

Rolling back the application is straightforward:

```bash
helm rollback leader-election 1
```

### CRD Rollback Issue and Resolution

When rolling back the CRD

```bash
helm rollback leader-election-crd 1
```

You will encounter the following error:

```console
Error: cannot patch "leaderelections.sample.operator.javaoperatorsdk.io" with kind CustomResourceDefinition: 
CustomResourceDefinition.apiextensions.k8s.io "leaderelections.sample.operator.javaoperatorsdk.io" is invalid: 
status.storedVersions[1]: Invalid value: "v1alpha2": must appear in spec.versions
```

To resolve this issue, you need to follow these steps before attempting the rollback:

1. Edit the CRD to update version information:

   ```bash
   kubectl edit crd leaderelections.sample.operator.javaoperatorsdk.io
   ```

2. In the editor:
   - Find the version you wish to downgrade from (v1alpha2) at index `spec.versions[1]` (as it's the second version in the array) and set:

     ```yaml
     spec.versions[1].served: false
     spec.versions[1].storage: false
     ```

   - Find the version you wish to roll back to (v1alpha1) at index `spec.versions[0]` (as it's the first version in the array) and set:

     ```yaml
     spec.versions[0].served: true
     spec.versions[0].storage: true
     ```

   - Save and exit the editor

   > **Note:** Alternatively, these two steps can be replaced by applying a pre-configured CRD yaml file:
   >
   > ```bash
   > kubectl apply -f pre-rollback-crd.yaml
   > ```

3. Patch the CRD to remove the v1alpha2 version from storedVersions:

   ```bash
   kubectl patch --subresource=status --type=json crd leaderelections.sample.operator.javaoperatorsdk.io -p '[{ "op": "replace", "path": "/status/storedVersions", "value": ["v1alpha1"]}]'
   ```

4. After these steps, you can successfully roll back the CRD:

   ```bash
   helm rollback leader-election-crd 1
   ```

> **Note:** This rollback procedure is based on the solution documented in [this HackMD article](https://hackmd.io/@6K3Bd7JbRlW_lllu-R4Xww/rk1KuN_JT).

## Troubleshooting

- If you encounter issues during installation or upgrades, check the Kubernetes and Helm logs.
- Verify the CRD status with `kubectl get crd leaderelections.sample.operator.javaoperatorsdk.io -o yaml`.
- For persistent issues, try reinstalling the CRD and application from scratch.

## Notes

The rollback issue occurs because Kubernetes validates that all versions in `status.storedVersions` must appear in `spec.versions`. When rolling back from v1alpha2 to v1alpha1, the stored v1alpha2 version is no longer present in the spec, causing the validation error. The manual editing and patching steps ensure proper version alignment before the rollback is attempted.