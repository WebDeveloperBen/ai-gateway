<script setup lang="ts">
import type { MetricSpec } from "~/components/charts/Card.vue"
import { rows } from "~/components/charts/demo-data"

// --- date helpers (unchanged) ---
const AUS_TZ = "Australia/Sydney"
const parseIsoDateOnlyToUtcMs = (iso: string) => {
  const [y, m, d] = iso.split("-").map(Number)
  return Date.UTC(y ?? 1, (m ?? 1) - 1, d ?? 1)
}
const formatAuDate = (msUtc: number) =>
  new Intl.DateTimeFormat("en-AU", {
    timeZone: AUS_TZ,
    year: "numeric",
    month: "short",
    day: "2-digit"
  }).format(msUtc)

const xLabels = computed(() => rows.value.map((r) => formatAuDate(parseIsoDateOnlyToUtcMs(r.date))))

// Assume a settings store exists; fall back if not.
const currency = (settings as any)?.value?.currency ?? "USD"

// Single source of truth for what to render:
const metrics = computed<MetricSpec[]>(() => [
  {
    key: "costUsd",
    title: "Cost",
    seriesName: "Cost (USD)",
    type: "line",
    prefix: "$",
    rounding: 2,
    compactAxis: true,
    valueMap: (r) => Number(r.costUsd.toFixed(2)),
    describe: (v) => `${v.toFixed(2)} ${currency}`
  },
  {
    key: "tokens",
    title: "Tokens",
    seriesName: "Tokens",
    type: "line",
    suffix: " tokens",
    compactAxis: true,
    describe: (v) => `${v.toLocaleString()} Used`
  },
  {
    key: "uniqueUsers",
    title: "Unique Users",
    seriesName: "Unique Users",
    type: "bar",
    suffix: " users",
    compactAxis: true,
    describe: (v) => v.toLocaleString()
  },
  {
    key: "requests",
    title: "Requests",
    seriesName: "Requests",
    type: "bar",
    suffix: " req",
    compactAxis: true,
    describe: (v) => v.toLocaleString()
  }
])
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <ChartsCard v-for="m in metrics" :key="m.key" :rows="rows" :labels="xLabels" :spec="m" />
  </div>
</template>
