# Delivery pipeline for the AICR provider

Registered as `pequod/aicr-provider-delivery/main`.

```
build ──► test ──► preview ──┬──► deploy    (manual + approval)
                             └──► teardown  (manual)
```

| Stage | Member | What | Cost | Gate |
|---|---|---|---|---|
| `build` | Job | `make provider` | free | — |
| `test` | 2 Jobs | `make test`, `make lint` in parallel | free | — |
| `preview` | Job | builds the Node SDK, wires an example to it, `pulumi preview --diff` | free | — |
| `deploy` | Job | `pulumi up` on the example | **~$196/hr** | `allOf(manual, approval)` |
| `teardown` | Job | `pulumi destroy` | free | `manual` |

`test` is two jobs so a `go vet` failure doesn't mask a unit-test failure.
`teardown` hangs off `preview`, not `deploy`, so it stays reachable when a
deploy dies half way — exactly when a p5 fleet is up and billing.

## Why this pipeline resolves its source

This repository is **public**. Delivery resolves public repos anonymously, with
no GitHub App installation involved. Pipelines pointed at internal repos fail
*silently* — no CLI error, no `Registered.` line — and the console misreports it
as a missing branch.

## Why it builds the provider instead of installing it

Every example depends on pulumi-org package names that were never published:

| Package | Registry | Status |
|---|---|---|
| `@pulumi/nvidia-aicr` | npm | **not published** |
| `pulumi-nvidia-aicr` | PyPI | **not published** |
| `@pulumi-labs/nvidia-aicr` | npm | 0.3.2 |
| `pulumi-labs-nvidia-aicr` | PyPI | 0.3.2 |

`npm install` in any example 404s today. The provider source is in this same
checkout, so the jobs run `make provider` and `make build_nodejs_sdk` and
install the result over whatever npm resolved.

That is also why `deploy` is a **Job and not a `Stacks` member**: Delivery
installs a Stacks member's dependencies *before* `preRunCommands` run, and those
dependencies are the unpublished packages. A Job keeps ordering under our
control, at the cost of the console's stack linkage, config injection and
artifacts. Once the pulumi-org packages ship, delete the build steps in
`jobs/_example-setup.sh` and convert `deploy` to a `Stacks` member.

## Starting a run

Registration keys the pipeline to your **local HEAD**. A run starts when that
commit becomes the head of the configured branch — registering is not the
trigger. Current state:

```
Registered. Commit 8633aaa is not the head of branch main (the head is 3276738).
Push 8633aaa to start a run for it.
```

So this needs to land on `main` (PR is the right shape for a shared repo). On
merge the run starts by itself; there is no button.

## Unverified

- **Go in the image.** `make provider` needs Go, `go.mod` asks for 1.27, and the
  jobs set `GOTOOLCHAIN=auto` so a present-but-older toolchain fetches it. The
  image is pinned to `pulumi/pulumi:3.241.0` (matching `mise.toml`), but I could
  not confirm it ships Go. If `build` fails with "go: not found", that is why —
  pick an image that has both Go and the Pulumi CLI.
- **Nothing has run yet.** Every stage above is reasoned from the Makefile and
  the Delivery schema, not observed. The first run is the real test.

## Credentials

Jobs open the ESC environment themselves —
`pulumi env run pequod/cloud-access/aws-access -- …` — so each gets a fresh OIDC
session rather than credentials frozen at registration time. pequod's role is in
account 052848974346, us-west-2.

## Before anyone clicks deploy

`examples/aws-eks-inference-ts` hardcodes `p5.48xlarge`; the only knob is
`nodeCount`. Two nodes is ~$196/hr. The gate states the number *and* requires a
named approver — `approval()` has no description field and `manual()` has no
approver list, so `allOf` of both is the only way to get each. pequod's
us-west-2 sandbox almost certainly has no p5 quota, so expect a capacity failure
before a bill. `build`, `test` and `preview` are free and are where the value is
until the packages ship.
