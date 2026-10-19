---
title: NVIDIA AI Cluster Runtime (AICR) Provider
meta_desc: Deploy validated NVIDIA AI Cluster Runtime (AICR) recipes on Kubernetes via a single Pulumi component.
layout: package
---

A Pulumi component provider that installs NVIDIA's GPU software stack on
Kubernetes. The `ClusterStack` component takes an accelerator, intent, and
target platform, and installs the corresponding Helm charts in dependency
order.

## Installation

The AICR provider is available as a package in all Pulumi-supported
languages:

* JavaScript/TypeScript: [`@pulumi/nvidia-aicr`](https://www.npmjs.com/package/@pulumi/nvidia-aicr)
* Python: [`pulumi-nvidia-aicr`](https://pypi.org/project/pulumi-nvidia-aicr/)
* Go: [`github.com/pulumi/pulumi-nvidia-aicr/sdk/go/nvidiaaicr`](https://github.com/pulumi/pulumi-nvidia-aicr)
* .NET: [`Pulumi.NvidiaAicr`](https://www.nuget.org/packages/Pulumi.NvidiaAicr)
* Java: [`com.pulumi:nvidia-aicr`](https://central.sonatype.com/artifact/com.pulumi/nvidia-aicr)

## Example Usage

The example below deploys an AICR recipe against an existing Kubernetes
cluster (here, a local `kind` cluster suitable for development; replace
`service` and provide a real `kubeconfig` for a cloud target).

{{< chooser language "typescript,python,go,csharp,java,yaml" >}}
{{% choosable language typescript %}}
```typescript
import * as aicr from "@pulumi/nvidia-aicr";

const stack = new aicr.ClusterStack("aicr", {
    accelerator: "h100",
    service: "kind",
    intent: "inference",
});
```
{{% /choosable %}}
{{% choosable language python %}}
```python
import pulumi_nvidia_aicr as aicr

aicr.ClusterStack(
    "aicr",
    accelerator="h100",
    service="kind",
    intent="inference",
)
```
{{% /choosable %}}
{{% choosable language go %}}
```go
package main

import (
    aicr "github.com/pulumi/pulumi-nvidia-aicr/sdk/go/nvidiaaicr"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := aicr.NewClusterStack(ctx, "aicr", &aicr.ClusterStackArgs{
            Accelerator: "h100",
            Service:     "kind",
            Intent:      "inference",
        })
        return err
    })
}
```
{{% /choosable %}}
{{% choosable language csharp %}}
```csharp
using Pulumi.NvidiaAicr;

return await Pulumi.Deployment.RunAsync(() =>
{
    var stack = new ClusterStack("aicr", new ClusterStackArgs
    {
        Accelerator = "h100",
        Service     = "kind",
        Intent      = "inference",
    });
});
```
{{% /choosable %}}
{{% choosable language java %}}
```java
import com.pulumi.Pulumi;
import com.pulumi.nvidiaaicr.ClusterStack;
import com.pulumi.nvidiaaicr.ClusterStackArgs;

public final class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            new ClusterStack("aicr", ClusterStackArgs.builder()
                .accelerator("h100")
                .service("kind")
                .intent("inference")
                .build());
        });
    }
}
```
{{% /choosable %}}
{{% choosable language yaml %}}
```yaml
resources:
  aicr:
    type: nvidia-aicr:index:ClusterStack
    properties:
      accelerator: h100
      service: kind
      intent: inference
```
{{% /choosable %}}
{{< /chooser >}}

For cloud-targeted scenarios (EKS, AKS, GKE, OKE, CoreWeave), see the
runnable per-language programs under
[`examples/`](https://github.com/pulumi/pulumi-nvidia-aicr/tree/main/examples)
in the provider repository.

## Configuration

`ClusterStack` is a Pulumi component resource; the resolved set of Helm
components is determined by the three criteria fields (`accelerator`,
`service`, `intent`) plus the recipe overlay system. See the API Reference
(left navigation) for the full input surface, including
`componentOverrides` for per-component Helm value customization.
