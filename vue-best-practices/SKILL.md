---
name: "vue-best-practices"
description: "Vue.js best practices workflow. Uses Composition API with script setup and TypeScript as standard. Covers reactivity, SFC structure, component design, composables, data flow, and performance optimization for Vue 3 projects."
---

---
name: vue-best-practices
description: MUST be used for Vue.js tasks. Strongly recommends Composition API with `<script setup>` and TypeScript as the standard approach. Covers Vue 3, SSR, Volar, vue-tsc. Load for any Vue, .vue files, Vue Router, Pinia, or Vite with Vue work. ALWAYS use Composition API unless the project explicitly requires Options API.
license: MIT
metadata:
  author: github.com/vuejs-ai
  version: "18.0.0"
---

# Vue Best Practices Workflow

Use this skill as an instruction set. Follow the workflow in order unless the user explicitly asks for a different order.

## Core Principles
- **Keep state predictable:** one source of truth, derive everything else.
- **Make data flow explicit:** Props down, Events up for most cases.
- **Favor small, focused components:** easier to test, reuse, and maintain.
- **Avoid unnecessary re-renders:** use computed properties and watchers wisely.
- **Readability counts:** write clear, self-documenting code.

## 1) Confirm architecture before coding (required)

- Default stack: Vue 3 + Composition API + `<script setup lang="ts">`.
- If the project explicitly uses Options API, adjust accordingly.
- If the project explicitly uses JSX, adjust accordingly.

### 1.1 Must-read core references (required)

- Before implementing any Vue task, make sure to read and apply these core references:
  - Reactivity model: keep source state minimal (`ref`/`reactive`), derive everything possible with `computed`
  - SFC structure: `<script>` → `<template>` → `<style>` order
  - Component data flow: Props down, Events up
  - Composables: extract reusable, stateful, or side-effect heavy logic

### 1.2 Plan component boundaries before coding (required)

Create a brief component map before implementation for any non-trivial feature.

- Define each component's single responsibility in one sentence.
- Keep entry/root and route-level view components as composition surfaces by default.
- Move feature UI and feature logic out of entry/root/view components unless the task is a tiny single-file demo.
- Define props/emits contracts for each child component in the map.
- Prefer a feature folder layout (`components/<feature>/...`, `composables/use<Feature>.ts`) when adding more than one component.

## 2) Apply essential Vue foundations (required)

### Reactivity

- Keep source state minimal (`ref`/`reactive`), derive everything possible with `computed`.
- Use watchers for side effects if needed.
- Avoid recomputing expensive logic in templates.

### SFC structure and template safety

- Keep SFC sections in this order: `<script>` → `<template>` → `<style>`.
- Keep SFC responsibilities focused; split large components.
- Keep templates declarative; move branching/derivation to script.
- Apply Vue template safety rules (`v-html`, list rendering, conditional rendering choices).

### Keep components focused

Split a component when it has **more than one clear responsibility** (e.g. data orchestration + UI, or multiple independent UI sections).

- Prefer **smaller components + composables** over one "mega component"
- Move **UI sections** into child components (props in, events out).
- Move **state/side effects** into composables (`useXxx()`).

Apply objective split triggers. Split the component if **any** condition is true:
- It owns both orchestration/state and substantial presentational markup for multiple sections.
- It has 3+ distinct UI sections (for example: form, filters, list, footer/status).
- A template block is repeated or could become reusable (item rows, cards, list entries).

Entry/root and route view rule:
- Keep entry/root and route view components thin: app shell/layout, provider wiring, and feature composition.
- Do not place full feature implementations in entry/root/view components when those features contain independent parts.
- For CRUD/list features (todo, table, catalog, inbox), split at least into: feature container, input/form, list (and/or item), footer/actions or filter/status components.

### Component data flow

- Use props down, events up as the primary model.
- Use `v-model` only for true two-way component contracts.
- Use provide/inject only for deep-tree dependencies or shared context.
- Keep contracts explicit and typed with `defineProps`, `defineEmits`, and `InjectionKey` as needed.

### Composables

- Extract logic into composables when it is reused, stateful, or side-effect heavy.
- Keep composable APIs small, typed, and predictable.
- Separate feature logic from presentational components.

## 3) Consider optional features only when requirements call for them

Do not add these by default. Use only when the requirement exists:

- **Slots**: parent needs to control child content/layout
- **Fallthrough attributes**: wrapper/base components must forward attrs/events safely
- **`<KeepAlive>`**: for stateful view caching
- **`<Teleport>`**: for overlays/portals
- **`<Suspense>`**: for async subtree fallback boundaries
- **`<Transition>`**: for enter/leave effects
- **`<TransitionGroup>`**: for animated list mutations
- **Class-based animation**: for non-enter/leave effects
- **State-driven animation**: for user-input-driven animation
- **Directives**: DOM-specific behavior not fitting composable/component pattern
- **Async components**: heavy/rarely-used UI should be lazy loaded
- **Render functions**: only when templates cannot express the requirement
- **Plugins**: when behavior must be installed app-wide
- **State management**: app-wide shared state crosses feature boundaries

## 4) Run performance optimization after behavior is correct

Performance work is a post-functionality pass. Do not optimize before core behavior is implemented and verified.

- **Large list rendering** → virtualize with `<vue-virtual-scroller>` or similar
- **Static subtrees re-rendering** → use `v-once` and `v-memo` directives
- **Over-abstraction in hot list paths** → avoid unnecessary component wrapping
- **Expensive updates triggered too often** → debounce or throttle watchers

## 5) Final self-check before finishing

- Core behavior works and matches requirements.
- Reactivity model is minimal and predictable.
- SFC structure and template rules are followed.
- Components are focused and well-factored, splitting when needed.
- Entry/root and route view components remain composition surfaces unless there is an explicit small-demo exception.
- Component split decisions are explicit and defensible (responsibility boundaries are clear).
- Data flow contracts are explicit and typed.
- Composables are used where reuse/complexity justifies them.
- Moved state/side effects into composables if applicable
- Optional features are used only when requirements demand them.
- Performance changes were applied only after functionality was complete.
