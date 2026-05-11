---
name: "microsoft-foundry"
description: "Deploy, evaluate, and manage Microsoft Foundry AI agents end-to-end: Docker build, ACR push, agent create/deploy/invoke, batch eval, continuous eval, prompt optimizer, dataset curation from traces, quota management, RBAC."
---

---
name: microsoft-foundry
description: "Deploy, evaluate, and manage Foundry agents end-to-end: Docker build, ACR push, hosted/prompt agent create, container start, batch eval, continuous eval, prompt optimizer workflows, agent.yaml, dataset curation from traces. USE FOR: deploy agent to Foundry, hosted agent, create agent, invoke agent, evaluate agent, run batch eval, continuous eval, continuous monitoring, continuous eval status, optimize prompt, improve prompt, prompt optimizer, optimize agent instructions, improve agent instructions, optimize system prompt, deploy model, Foundry project, RBAC, role assignment, permissions, quota, capacity, region, troubleshoot agent, deployment failure, create dataset from traces, dataset versioning, eval trending, create AI Services, Cognitive Services, create Foundry resource, provision resource, knowledge index, agent monitoring, customize deployment, onboard, availability. DO NOT USE FOR: Azure Functions, App Service, general Azure deploy (use azure-deploy), general Azure prep (use azure-prepare)."
license: MIT
metadata:
  author: Microsoft
  version: "1.1.9"
---

# Microsoft Foundry Skill

This skill helps developers work with Microsoft Foundry resources, covering model discovery and deployment, complete dev lifecycle of AI agent, evaluation workflows, and troubleshooting.

## Pre-Execution Requirements

> **MANDATORY: Before executing ANY workflow, you MUST first call the Azure MCP `foundry` tool and inspect the available Foundry MCP tools and related parameters.** Treat this initial `foundry` call as a discovery/help step. For this skill, Azure MCP `foundry` is the required entry point for Foundry-related MCP operations.

## Sub-Skills

> **MANDATORY: Before executing ANY workflow-specific steps, you MUST read the corresponding sub-skill document.** Do not call workflow-specific MCP tools for a workflow without reading its skill document. This applies even if you already know the MCP tool parameters — the skill document contains required workflow steps, pre-checks, and validation logic that must be followed. This rule applies on every new user message that triggers a different workflow, even if the skill is already loaded.

This skill includes specialized sub-skills for specific workflows. **Use these instead of the main skill when they match your task:**

| Sub-Skill | When to Use | Reference |
|-----------|-------------|-----------|
| **deploy** | Containerize, build, push to ACR, create/update/clone agent deployments | [deploy](foundry-agent/deploy/deploy.md) |
| **invoke** | Send messages to an agent, single or multi-turn conversations | [invoke](foundry-agent/invoke/invoke.md) |
| **observe** | Evaluate agent quality, run batch evals, analyze failures, optimize prompts, improve agent instructions, compare versions, set up CI/CD monitoring, and enable continuous production evaluation | [observe](foundry-agent/observe/observe.md) |
| **trace** | Query traces, analyze latency/failures, correlate eval results to specific responses via App Insights `customEvents` | [trace](foundry-agent/trace/trace.md) |
| **troubleshoot** | View hosted agent logs, query telemetry, diagnose failures | [troubleshoot](foundry-agent/troubleshoot/troubleshoot.md) |
| **create** | Create new hosted agent applications. Supports Microsoft Agent Framework, LangGraph, or custom frameworks in Python or C#, across `responses` or `invocations` protocols. | [create](foundry-agent/create/create.md) |
| **eval-datasets** | Harvest production traces into evaluation datasets, manage dataset versions and splits, track evaluation metrics over time, detect regressions, and maintain full lineage from trace to deployment. | [eval-datasets](foundry-agent/eval-datasets/eval-datasets.md) |
| **project/create** | Creating a new Azure AI Foundry project for hosting agents and models. Use when onboarding to Foundry or setting up new infrastructure. | [project/create/create-foundry-project.md](project/create/create-foundry-project.md) |
| **resource/create** | Creating Azure AI Services multi-service resource (Foundry resource) using Azure CLI. Use when manually provisioning AI Services resources with granular control. | [resource/create/create-foundry-resource.md](resource/create/create-foundry-resource.md) |
| **private-network** | Answer questions about Foundry network isolation **and** deploy Foundry with VNet isolation (BYO VNet, Managed VNet, hybrid). | [resource/private-network/private-network.md](resource/private-network/private-network.md) |
| **models/deploy-model** | Unified model deployment with intelligent routing. Handles quick preset deployments, fully customized deployments, and capacity discovery. | [models/deploy-model/SKILL.md](models/deploy-model/SKILL.md) |
| **quota** | Managing quotas and capacity for Microsoft Foundry resources. | [quota/quota.md](quota/quota.md) |
| **rbac** | Managing RBAC permissions, role assignments, managed identities, and service principals for Microsoft Foundry resources. | [rbac/rbac.md](rbac/rbac.md) |

> 💡 **Tip:** For a complete onboarding flow: `project/create` (public) or `private-network` (VNet isolation) → `models/deploy-model` → agent workflows (`create` → `deploy` → `invoke`).

## Infrastructure Lifecycle

| User Intent | Workflow |
|-------------|---------|
| "Create Foundry" / "Set up Foundry" (ambiguous) | Use `AskUserQuestion`: (a) just an AI Services resource, (b) a project with public access, or (c) a project with network isolation? |
| Set up Foundry with VNet isolation | private-network |
| Create a Foundry project (public) | project/create |
| Create a bare Foundry resource | resource/create |

## Agent Development Lifecycle

| User Intent | Workflow (read in order) |
|-------------|------------------------|
| Create a new agent from scratch | create → deploy → invoke |
| Deploy an agent (code already exists) | deploy → invoke |
| Update/redeploy an agent after code changes | deploy → invoke |
| Invoke/test/chat with an agent | invoke |
| Optimize / improve agent prompt or instructions | observe (Step 4: Optimize) |
| Evaluate and optimize agent (full loop) | observe |
| Enable continuous evaluation monitoring | observe (Step 6: CI/CD & Monitoring) |
| Troubleshoot an agent issue | invoke → troubleshoot |
| Fix a broken agent (troubleshoot + redeploy) | invoke → troubleshoot → apply fixes → deploy → invoke |

## Agent: .foundry Workspace Standard

Every agent source folder should keep Foundry-specific state under `.foundry/`:

```text
<agent-root>/
  .foundry/
    agent-metadata.yaml
    agent-metadata.prod.yaml
    datasets/
    evaluators/
    results/
```

- `agent-metadata.yaml` is the preferred local/dev metadata file.
- `datasets/` and `evaluators/` are local cache folders.

## Tool Usage Conventions

- Use the `AskUserQuestion` tool whenever collecting information from the user
- Prefer Azure MCP tools over direct CLI commands when available

## SDK Quick Reference

- See references for Python SDK documentation.
