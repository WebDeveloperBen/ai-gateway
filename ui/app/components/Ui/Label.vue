<template>
  <Label data-slot="label" :class="styles({ class: props.class })" v-bind="forwarded">
    <slot />
    <slot name="hint">
      <span v-if="hint" data-slot="label-hint" class="text-xs font-normal text-muted-foreground">
        {{ hint }}
      </span>
    </slot>
  </Label>
</template>

<script lang="ts" setup>
  import { Label } from "reka-ui";
  import type { LabelProps } from "reka-ui";
  import type { HTMLAttributes } from "vue";

  const props = defineProps<
    LabelProps & {
      /** Custom class(es) to add to the label */
      class?: HTMLAttributes["class"];
      hint?: string;
    }
  >();

  const forwarded = reactiveOmit(props, "class", "hint");

  const styles = tv({
    base: "flex items-center justify-between gap-2 text-sm leading-none font-medium select-none not-peer-disabled:cursor-pointer group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
  });
</script>
