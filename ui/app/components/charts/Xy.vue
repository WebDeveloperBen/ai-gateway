<script setup lang="ts">
import {
  VueDataUi,
  type VueUiXyDatasetItem,
  type VueUiXyConfig,
} from "vue-data-ui";

type ChartType = "line" | "bar" | "plot";

const props = defineProps<{
  labels: string[];
  series: number[];
  seriesName?: string;
  type?: ChartType;
  prefix?: string;
  suffix?: string;
  rounding?: number;
  compactAxis?: boolean;
}>();

const fmtCompact = (n: number, maxFrac = 1) =>
  new Intl.NumberFormat("en", {
    notation: "compact",
    maximumFractionDigits: maxFrac,
  })
    .format(n)
    .replace("K", "k");

function asDataset(
  values: number[],
  name: string,
  type: ChartType = "line",
  dataLabels = false,
): VueUiXyDatasetItem[] {
  return [{ name, series: values, type, dataLabels }];
}

function xyConfig(
  labels: string[],
  opts?: {
    rounding?: number;
    prefix?: string;
    suffix?: string;
    compactAxis?: boolean;
  },
): VueUiXyConfig {
  const rounding = opts?.rounding ?? 0;
  const compactAxis = opts?.compactAxis ?? true;
  return {
    responsive: true,
    responsiveProportionalSizing: true,
    customPalette: [
      "var(--chart-1)",
      "var(--chart-2)",
      "var(--chart-3)",
      "var(--chart-4)",
      "var(--chart-5)",
    ],
    chart: {
      fontFamily: "var(--font-sans)",
      backgroundColor: "transparent",
      color: "var(--foreground)",
      grid: {
        stroke: "var(--border)",
        labels: {
          color: "var(--muted-foreground)",
          xAxisLabels: {
            values: labels,
            show: true,
            rotation: 90,
            autoRotate: { enable: true },
            color: "var(--muted-foreground)",
          },
          xAxis: {
            showBaseline: false,
            showCrosshairs: false,
            crosshairsAlwaysAtZero: true,
          },
          yAxis: {
            rounding,
            formatter: (params: any) => {
              const v =
                typeof params?.value === "number"
                  ? params.value
                  : Number(params?.value);
              if (!Number.isFinite(v)) return "–";
              return compactAxis ? fmtCompact(v) : v.toLocaleString();
            },
          },
        },
      },
      legend: {
        show: false,
        position: "bottom",
        color: "var(--muted-foreground)",
      },
      tooltip: { show: true, showTimeLabel: true, showPercentage: false },
      userOptions: { show: false },
      zoom: { show: false },
      highlighter: { useLine: true, opacity: 3 },
      labels: { prefix: opts?.prefix ?? "", suffix: opts?.suffix ?? "" },
    },
    line: { area: { useGradient: true, opacity: 30 } },
    bar: { useGradient: true },
  } as VueUiXyConfig;
}

const dataset = computed(() =>
  asDataset(props.series, props.seriesName ?? "Series", props.type ?? "line"),
);

const config = computed(() =>
  xyConfig(props.labels, {
    prefix: props.prefix,
    suffix: props.suffix,
    rounding: props.rounding,
    compactAxis: props.compactAxis ?? true,
  }),
);
</script>

<template>
  <VueDataUi
    component="VueUiXy"
    :dataset="dataset"
    :config="config"
    class="w-full h-full"
  />
</template>
