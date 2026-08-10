# Provider Vultr

`provider-vultr` is a [Crossplane](https://crossplane.io/) provider for
[Vultr](https://www.vultr.com/), built with
[Upjet](https://github.com/crossplane/upjet) and backed by the
[`vultr/vultr`](https://github.com/vultr/terraform-provider-vultr)
Terraform provider.

It exposes XRM-conformant managed resources that let you manage Vultr
infrastructure — instances, Kubernetes clusters (VKE), block storage,
load balancers, DNS, and more — directly from Kubernetes or Upbound.
Every resource is available in two flavors: cluster-scoped
(`*.vultr.upbound.io`) and namespaced (`*.vultr.m.upbound.io`).

## Authentication

The provider authenticates with a Vultr API key read from the referenced
`Secret`. Generate one in the [Vultr customer portal](https://my.vultr.com/settings/#settingsapi)
and store it as JSON:

```json
{
  "api_key": "<Vultr API key>"
}
```

## Getting Started

### 1. Install the provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-vultr
spec:
  package: xpkg.upbound.io/upbound/provider-vultr:v1.0.0
```

### 2. Create a credentials Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vultr-creds
  namespace: upbound-system
type: Opaque
stringData:
  creds: |
    {
      "api_key": "<Vultr API key>"
    }
```

### 3. Create a ProviderConfig

```yaml
apiVersion: vultr.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: vultr-creds
      namespace: upbound-system
      key: creds
```

### 4. Create a managed resource

More examples for each resource are available under [`examples/cluster/`](examples/cluster/)
and [`examples/namespaced/`](examples/namespaced/).

## ProviderConfig fields

| Field | Required | Description |
|---|---|---|
| `spec.credentials.source` | Yes | One of `Secret`, `InjectedIdentity`, `Environment`, `Filesystem` |
| `spec.credentials.secretRef` | When source=Secret | Reference to the credentials Secret |
| `spec.rateLimit` | No | Speed of API calls in milliseconds, to work within the Vultr rate limit |
| `spec.retryLimit` | No | Maximum number of retries for a failed API call |
| `spec.reconciliationPolicy` | No | Rate-limiting policy for reconciliation |

## Developing

### Code generation

```console
make generate
```

This runs the Upjet code generator against the pinned `vultr/vultr`
Terraform provider schema and docs, and writes the generated APIs,
controllers, and CRDs for both the cluster-scoped and namespaced variants.

### Run locally against a cluster

```console
make run
```

### Run end-to-end tests

```console
make e2e
```

## Reporting issues

Please open an [issue](https://github.com/upbound/provider-upjet-vultr/issues)
for bug reports, feature requests, or questions.
