<script setup lang="ts">
type Row = {
  date: string;
  costUsd: number;
  tokens: number;
  uniqueUsers: number;
  requests: number;
};

type ChartType = "line" | "bar" | "plot";
type ChartableKey = Exclude<keyof Row, "date">;

export type MetricSpec<K extends ChartableKey = ChartableKey> = {
  key: K;
  title: string;
  seriesName?: string;
  type?: ChartType;
  prefix?: string;
  suffix?: string;
  rounding?: number;
  compactAxis?: boolean;
  valueMap?: (r: Row) => number;
  describe?: (v: number) => string;
};

const props = defineProps<{
  rows: Row[];
  labels: string[];
  spec: MetricSpec;
}>();

const series = computed<number[]>(() =>
  props.rows.map((r) =>
    props.spec.valueMap
      ? props.spec.valueMap(r)
      : (r[props.spec.key] as number),
  ),
);

const currentText = computed(() => {
  const v = series.value.at(-1);
  if (v == null) return "-";
  return props.spec.describe ? props.spec.describe(v) : v.toLocaleString();
});
</script>

<template>
  <UiCard>
    <UiCardHeader>
      <span class="font-semibold">{{ spec.title }}</span>
      <UiCardDescription>{{ currentText }}</UiCardDescription>
    </UiCardHeader>

    <UiCardContent>
      <ClientOnly>
        <div class="w-full h-[260px] sm:min-h-[290px] 3xl:h-[450px]">
          <ChartsXy
            :labels="labels"
            :series="series"
            :seriesName="spec.seriesName"
            :type="spec.type"
            :prefix="spec.prefix"
            :suffix="spec.suffix"
            :rounding="spec.rounding"
            :compactAxis="spec.compactAxis"
          />
        </div>
      </ClientOnly>
    </UiCardContent>
  </UiCard>
</template>
