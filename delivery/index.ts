// Delivery pipeline for the Pulumi NVIDIA AICR provider, registered in pequod.
//
//   build ──► test ──► preview ──┬──► deploy    (manual + approval)
//                                └──► teardown  (manual)
//
// The source is this repository, which is public. Delivery resolves public
// repos anonymously, with no GitHub App installation involved -- which is why
// this pipeline resolves where one pointed at an internal repo silently does
// not.
//
// The pipeline builds the provider rather than installing it because every
// example here depends on the pulumi-org package names (@pulumi/nvidia-aicr,
// pulumi-nvidia-aicr, Pulumi.NvidiaAicr) and none of them are published -- npm
// and PyPI carry only the pulumi-labs names, at 0.3.2. The provider source is
// in this same checkout, so `make provider` is both available and necessary.

import * as pulumi from "@pulumi/pulumi";
import * as delivery from "@pulumi/delivery";

const cfg = new pulumi.Config();

const repository = cfg.get("repository") ?? "pulumi/pulumi-nvidia-aicr";
const branch = cfg.get("branch") ?? "delivery-pipeline";
const approvers = cfg.getObject<string[]>("approvers") ?? [];

const example = cfg.get("example") ?? "examples/aws-eks-inference-ts";
const awsEnvironment = cfg.get("awsEnvironment") ?? "pequod/cloud-access/aws-access";
const targetStack = cfg.get("targetStack") ?? "pequod/aws-eks-h100-inference/dev";

// `make provider` needs Go and the SDK build needs Node; the default Delivery
// image carries only the Pulumi CLI. 3.241.0 matches the pulumi version
// mise.toml pins, so SDK regeneration here matches CI.
const image = { reference: cfg.get("image") ?? "pulumi/pulumi:3.241.0" };

// go.mod asks for Go 1.27. GOTOOLCHAIN=auto lets the toolchain in the image
// fetch it rather than failing on a version mismatch.
const goEnv = { GOTOOLCHAIN: "auto", GOFLAGS: "-mod=mod" };

const exampleEnv = {
    ...goEnv,
    EXAMPLE_DIR: example,
    TARGET_STACK: targetStack,
    AWS_ENVIRONMENT: awsEnvironment,
};

const src = new delivery.GitHubSource("provider", {
    repository,
    branch,
    // Docs and assets do not change the provider, so they do not start a run.
    trigger: { paths: ["provider/**", "sdk/**", "examples/**", "delivery/**", "Makefile"] },
});

const pipeline = new delivery.Pipeline("aicr-provider", {
    program: { source: src, directory: "delivery" },
});

// --- 1. Build --------------------------------------------------------------
// Everything downstream needs the binary, and this is the cheapest thing that
// can fail, so it runs first and alone.
const buildStage = new delivery.Stage("build", { pipeline });

new delivery.Job("provider", {
    stage: buildStage,
    source: src,
    directory: ".",
    image,
    environment: goEnv,
    command: "make provider && ls -l bin/",
});

// --- 2. Test ---------------------------------------------------------------
// Two jobs in one stage so they run in parallel and a vet failure does not hide
// a test failure. Both are plain Go, so the image needs nothing extra.
const testStage = new delivery.Stage("test", { pipeline, from: [buildStage] });

new delivery.Job("unit", {
    stage: testStage,
    source: src,
    directory: ".",
    image,
    environment: goEnv,
    command: "make test",
});

new delivery.Job("vet", {
    stage: testStage,
    source: src,
    directory: ".",
    image,
    environment: goEnv,
    command: "make lint",
});

// --- 3. Preview ------------------------------------------------------------
// Builds the Node SDK, points an example at it, and prints the plan. Creates
// nothing. This is the first point where the provider is exercised end to end:
// schema, SDK and resource registration against a real program.
const previewStage = new delivery.Stage("preview", { pipeline, from: [testStage] });

new delivery.Job("example", {
    stage: previewStage,
    source: src,
    directory: ".",
    image,
    environment: exampleEnv,
    timeout: "30m",
    command: "sh delivery/jobs/preview-example.sh",
});

// --- 4. Deploy -------------------------------------------------------------
// aws-eks-inference hardcodes p5.48xlarge -- ~$98/hr per node, two nodes. The
// gate is a manual acknowledgement carrying that number *and* a named approval,
// because approval() has no description field and manual() has no approver
// list. allOf gives both.
const deployGate = Promise.all([
    delivery.manual({
        description:
            "Deploy for real? examples/aws-eks-inference runs 2x p5.48xlarge at " +
            "~$98/hr each -- about $196/hr until teardown.",
    }),
    delivery.approval({ approvers }),
]).then((rules) => delivery.allOf({ rules }));

const deployStage = new delivery.Stage("deploy", {
    pipeline,
    from: [previewStage],
    when: deployGate,
});

// A Job, not a Stacks member: Delivery installs a Stacks member's dependencies
// before preRunCommands run, and those dependencies are the unpublished
// pulumi-org packages. A Job keeps ordering under our control, at the cost of
// the console's stack linkage and artifacts. Swap it once the packages ship.
new delivery.Job("deploy", {
    stage: deployStage,
    source: src,
    directory: ".",
    image,
    environment: exampleEnv,
    timeout: "90m",
    command: "sh delivery/jobs/deploy-example.sh",
});

// --- 5. Teardown -----------------------------------------------------------
// Hangs off preview rather than deploy so it stays reachable when a deploy dies
// half way -- which is exactly when a p5 fleet is up and billing.
const teardownStage = new delivery.Stage("teardown", {
    pipeline,
    from: [previewStage],
    when: delivery.manual({
        description: "Destroy the example stack. Safe to run even if deploy failed.",
    }),
});

new delivery.Job("teardown", {
    stage: teardownStage,
    source: src,
    directory: ".",
    image,
    environment: exampleEnv,
    timeout: "60m",
    command: "sh delivery/jobs/teardown-example.sh",
});

export const release = pipeline.releaseNumber;
export const commit = src.commit;
export const sourceBranch = branch;
