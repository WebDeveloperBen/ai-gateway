<script lang="ts">
import type { Ref, WatchOptions } from "vue"
import type { Cell, Header, RowData, TableMeta } from "@tanstack/table-core"
import type {
  CellContext,
  ColumnDef,
  ColumnFiltersOptions,
  ColumnFiltersState,
  ColumnOrderState,
  ColumnPinningOptions,
  ColumnPinningState,
  ColumnSizingInfoState,
  ColumnSizingOptions,
  ColumnSizingState,
  CoreOptions,
  ExpandedOptions,
  ExpandedState,
  FacetedOptions,
  GlobalFilterOptions,
  GroupingOptions,
  GroupingState,
  HeaderContext,
  PaginationOptions,
  PaginationState,
  Row,
  RowPinningOptions,
  RowPinningState,
  RowSelectionOptions,
  RowSelectionState,
  SortingOptions,
  SortingState,
  Updater,
  VisibilityOptions,
  VisibilityState
} from "@tanstack/vue-table"

declare module "@tanstack/table-core" {
  interface ColumnMeta<TData extends RowData, TValue> {
    class?: {
      th?: string | ((cell: Header<TData, TValue>) => string)
      td?: string | ((cell: Cell<TData, TValue>) => string)
    }
    style?: {
      th?: string | Record<string, string> | ((cell: Header<TData, TValue>) => string | Record<string, string>)
      td?: string | Record<string, string> | ((cell: Cell<TData, TValue>) => string | Record<string, string>)
    }
    colspan?: {
      td?: string | ((cell: Cell<TData, TValue>) => string)
    }
    rowspan?: {
      td?: string | ((cell: Cell<TData, TValue>) => string)
    }
  }

  interface TableMeta<TData> {
    class?: {
      tr?: string | ((row: Row<TData>) => string)
    }
    style?: {
      tr?: string | Record<string, string> | ((row: Row<TData>) => string | Record<string, string>)
    }
  }
}

export type TableRow<T> = Row<T>
export type TableData = RowData
export type TableColumn<T extends TableData, D = unknown> = ColumnDef<T, D>

export interface TableOptions<T extends TableData = TableData>
  extends Omit<
    CoreOptions<T>,
    "data" | "columns" | "getCoreRowModel" | "state" | "onStateChange" | "renderFallbackValue"
  > {
  state?: CoreOptions<T>["state"]
  onStateChange?: CoreOptions<T>["onStateChange"]
  renderFallbackValue?: CoreOptions<T>["renderFallbackValue"]
}

export interface TableProps<T extends TableData = TableData> extends TableOptions<T> {
  as?: any
  data?: T[]
  columns?: TableColumn<T>[]
  caption?: string
  meta?: TableMeta<T>
  empty?: string
  sticky?: boolean | "header" | "footer"
  loading?: boolean
  // keep props but loosen types to plain strings so we don't depend on app config types
  loadingColor?: string
  loadingAnimation?: string
  watchOptions?: WatchOptions
  globalFilterOptions?: Omit<GlobalFilterOptions<T>, "onGlobalFilterChange">
  columnFiltersOptions?: Omit<ColumnFiltersOptions<T>, "getFilteredRowModel" | "onColumnFiltersChange">
  columnPinningOptions?: Omit<ColumnPinningOptions, "onColumnPinningChange">
  columnSizingOptions?: Omit<ColumnSizingOptions, "onColumnSizingChange" | "onColumnSizingInfoChange">
  visibilityOptions?: Omit<VisibilityOptions, "onColumnVisibilityChange">
  sortingOptions?: Omit<SortingOptions<T>, "getSortedRowModel" | "onSortingChange">
  groupingOptions?: Omit<GroupingOptions, "onGroupingChange">
  expandedOptions?: Omit<ExpandedOptions<T>, "getExpandedRowModel" | "onExpandedChange">
  rowSelectionOptions?: Omit<RowSelectionOptions<T>, "onRowSelectionChange">
  rowPinningOptions?: Omit<RowPinningOptions<T>, "onRowPinningChange">
  paginationOptions?: Omit<PaginationOptions, "onPaginationChange">
  facetedOptions?: FacetedOptions<T>
  onSelect?: (row: TableRow<T>, e?: Event) => void
  onHover?: (e: Event, row: TableRow<T> | null) => void
  onContextmenu?: ((e: Event, row: TableRow<T>) => void) | Array<(e: Event, row: TableRow<T>) => void>
  class?: any
  // keep `ui` passthrough for user overrides; just treat it as a bag of classnames
  ui?: {
    root?: any
    base?: any
    caption?: any
    thead?: any
    tr?: any
    th?: any
    separator?: any
    tbody?: any
    td?: any
    loading?: any
    empty?: any
    tfoot?: any
  }
}

type DynamicHeaderSlots<T, K = keyof T> = Record<string, (props: HeaderContext<T, unknown>) => any> &
  Record<`${K extends string ? K : never}-header`, (props: HeaderContext<T, unknown>) => any>
type DynamicFooterSlots<T, K = keyof T> = Record<string, (props: HeaderContext<T, unknown>) => any> &
  Record<`${K extends string ? K : never}-footer`, (props: HeaderContext<T, unknown>) => any>
type DynamicCellSlots<T, K = keyof T> = Record<string, (props: CellContext<T, unknown>) => any> &
  Record<`${K extends string ? K : never}-cell`, (props: CellContext<T, unknown>) => any>

export type TableSlots<T extends TableData = TableData> = {
  expanded: (props: { row: Row<T> }) => any
  empty: (props?: {}) => any
  loading: (props?: {}) => any
  caption: (props?: {}) => any
  "body-top": (props?: {}) => any
  "body-bottom": (props?: {}) => any
} & DynamicHeaderSlots<T> &
  DynamicFooterSlots<T> &
  DynamicCellSlots<T>
</script>

<script setup lang="ts" generic="T extends TableData">
import { computed, ref, watch } from "vue"
import { Primitive } from "reka-ui"
import { upperFirst } from "scule"
import {
  FlexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  getExpandedRowModel,
  useVueTable
} from "@tanstack/vue-table"
import { reactiveOmit } from "@vueuse/core"

const props = withDefaults(defineProps<TableProps<T>>(), {
  watchOptions: () => ({
    deep: true
  })
})

// ❌ removed: const appConfig = useAppConfig() as Table["AppConfig"]

const data = ref(props.data ?? []) as Ref<T[]>
const columns = computed<TableColumn<T>[]>(
  () =>
    props.columns ??
    Object.keys(data.value[0] ?? {}).map((accessorKey: string) => ({ accessorKey, header: upperFirst(accessorKey) }))
)
const meta = computed(() => props.meta ?? {})

/** tiny class combiner (shadcn-style) */
function cn(...inputs: any[]) {
  return inputs.flat(Infinity).filter(Boolean).join(" ")
}

/** Tailwind/shadcn-flavored UI map mirroring your previous `ui.*` API */
const ui = computed(() => {
  const isStickyHeader = props.sticky === true || props.sticky === "header"
  const isStickyFooter = props.sticky === true || props.sticky === "footer"

  return {
    root: ({ class: extra }: any = {}) =>
      cn("relative w-full overflow-auto rounded-2xl border bg-card text-card-foreground shadow-sm", extra),
    base: ({ class: extra }: any = {}) => cn("w-full caption-bottom text-sm", extra),
    caption: ({ class: extra }: any = {}) => cn("mt-4 text-xs text-muted-foreground", extra),

    thead: ({ class: extra }: any = {}) =>
      cn(
        "text-sm",
        isStickyHeader && "sticky top-0 z-10 bg-muted/60 backdrop-blur supports-[backdrop-filter]:bg-muted/40",
        extra
      ),
    tfoot: ({ class: extra }: any = {}) =>
      cn(
        "text-sm",
        isStickyFooter && "sticky bottom-0 z-10 bg-muted/60 backdrop-blur supports-[backdrop-filter]:bg-muted/40",
        extra
      ),

    tr: ({ class: extra }: any = {}) =>
      cn(
        "border-b last:border-0 data-[selectable=true]:cursor-pointer",
        "hover:bg-muted/30 data-[selected=true]:bg-muted data-[expanded=true]:bg-muted/20",
        extra
      ),

    th: ({ pinned, class: extra }: any = {}) =>
      cn(
        "h-10 px-3 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0",
        pinned && "sticky left-0 bg-background shadow-[inset_-1px_0_0_theme(colors.border)]",
        extra
      ),

    td: ({ pinned, class: extra }: any = {}) =>
      cn(
        "p-3 align-middle [&:has([role=checkbox])]:pr-0",
        "whitespace-nowrap",
        pinned && "sticky left-0 bg-background shadow-[inset_-1px_0_0_theme(colors.border)]",
        extra
      ),

    separator: ({ class: extra }: any = {}) => cn("h-px bg-border", extra),
    tbody: ({ class: extra }: any = {}) => cn("", extra),

    loading: ({ class: extra }: any = {}) => cn("py-10 text-center text-muted-foreground", extra),
    empty: ({ class: extra }: any = {}) => cn("py-10 text-center text-muted-foreground", extra)
  }
})

const hasFooter = computed(() => {
  function hasFooterRecursive(columns: TableColumn<T>[]): boolean {
    for (const column of columns) {
      if ("footer" in column) {
        return true
      }
      if ("columns" in column && hasFooterRecursive(column.columns as TableColumn<T>[])) {
        return true
      }
    }
    return false
  }

  return hasFooterRecursive(columns.value)
})

const globalFilterState = defineModel<string>("globalFilter", { default: undefined })
const columnFiltersState = defineModel<ColumnFiltersState>("columnFilters", { default: [] })
const columnOrderState = defineModel<ColumnOrderState>("columnOrder", { default: [] })
const columnVisibilityState = defineModel<VisibilityState>("columnVisibility", { default: {} })
const columnPinningState = defineModel<ColumnPinningState>("columnPinning", { default: {} })
const columnSizingState = defineModel<ColumnSizingState>("columnSizing", { default: {} })
const columnSizingInfoState = defineModel<ColumnSizingInfoState>("columnSizingInfo", { default: {} })
const rowSelectionState = defineModel<RowSelectionState>("rowSelection", { default: {} })
const rowPinningState = defineModel<RowPinningState>("rowPinning", { default: {} })
const sortingState = defineModel<SortingState>("sorting", { default: [] })
const groupingState = defineModel<GroupingState>("grouping", { default: [] })
const expandedState = defineModel<ExpandedState>("expanded", { default: {} })
const paginationState = defineModel<PaginationState>("pagination", { default: {} })

const tableRef = ref<HTMLTableElement | null>(null)

const tableApi = useVueTable({
  ...reactiveOmit(
    props,
    "as",
    "data",
    "columns",
    "caption",
    "sticky",
    "loading",
    "loadingColor",
    "loadingAnimation",
    "class",
    "ui"
  ),
  data,
  get columns() {
    return columns.value
  },
  meta: meta.value,
  getCoreRowModel: getCoreRowModel(),
  ...(props.globalFilterOptions || {}),
  onGlobalFilterChange: (updaterOrValue) => valueUpdater(updaterOrValue, globalFilterState),
  ...(props.columnFiltersOptions || {}),
  getFilteredRowModel: getFilteredRowModel(),
  onColumnFiltersChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnFiltersState),
  onColumnOrderChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnOrderState),
  ...(props.visibilityOptions || {}),
  onColumnVisibilityChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnVisibilityState),
  ...(props.columnPinningOptions || {}),
  onColumnPinningChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnPinningState),
  ...(props.columnSizingOptions || {}),
  onColumnSizingChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnSizingState),
  onColumnSizingInfoChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnSizingInfoState),
  ...(props.rowSelectionOptions || {}),
  onRowSelectionChange: (updaterOrValue) => valueUpdater(updaterOrValue, rowSelectionState),
  ...(props.rowPinningOptions || {}),
  onRowPinningChange: (updaterOrValue) => valueUpdater(updaterOrValue, rowPinningState),
  ...(props.sortingOptions || {}),
  getSortedRowModel: getSortedRowModel(),
  onSortingChange: (updaterOrValue) => valueUpdater(updaterOrValue, sortingState),
  ...(props.groupingOptions || {}),
  onGroupingChange: (updaterOrValue) => valueUpdater(updaterOrValue, groupingState),
  ...(props.expandedOptions || {}),
  getExpandedRowModel: getExpandedRowModel(),
  onExpandedChange: (updaterOrValue) => valueUpdater(updaterOrValue, expandedState),
  ...(props.paginationOptions || {}),
  onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, paginationState),
  ...(props.facetedOptions || {}),
  state: {
    get globalFilter() {
      return globalFilterState.value
    },
    get columnFilters() {
      return columnFiltersState.value
    },
    get columnOrder() {
      return columnOrderState.value
    },
    get columnVisibility() {
      return columnVisibilityState.value
    },
    get columnPinning() {
      return columnPinningState.value
    },
    get expanded() {
      return expandedState.value
    },
    get rowSelection() {
      return rowSelectionState.value
    },
    get sorting() {
      return sortingState.value
    },
    get grouping() {
      return groupingState.value
    },
    get rowPinning() {
      return rowPinningState.value
    },
    get columnSizing() {
      return columnSizingState.value
    },
    get columnSizingInfo() {
      return columnSizingInfoState.value
    },
    get pagination() {
      return paginationState.value
    }
  }
})

function valueUpdater<T extends Updater<any>>(updaterOrValue: T, ref: Ref) {
  ref.value = typeof updaterOrValue === "function" ? updaterOrValue(ref.value) : updaterOrValue
}

function onRowSelect(e: Event, row: TableRow<T>) {
  if (!props.onSelect) {
    return
  }
  const target = e.target as HTMLElement
  const isInteractive = target.closest("button") || target.closest("a")
  if (isInteractive) {
    return
  }

  e.preventDefault()
  e.stopPropagation()

  // FIXME: `e` should be the first argument for consistency
  props.onSelect(row, e)
}

function onRowHover(e: Event, row: TableRow<T> | null) {
  if (!props.onHover) {
    return
  }

  props.onHover(e, row)
}

function onRowContextmenu(e: Event, row: TableRow<T>) {
  if (!props.onContextmenu) {
    return
  }

  if (Array.isArray(props.onContextmenu)) {
    props.onContextmenu.forEach((fn) => fn(e, row))
  } else {
    props.onContextmenu(e, row)
  }
}

function resolveValue<T, A = undefined>(prop: T | ((arg: A) => T), arg?: A): T | undefined {
  if (typeof prop === "function") {
    // @ts-expect-error: TS can't know if prop is a function here
    return prop(arg)
  }
  return prop
}

watch(
  () => props.data,
  () => {
    data.value = props.data ? [...props.data] : []
  },
  props.watchOptions
)

defineExpose({
  tableRef,
  tableApi
})
</script>

<template>
  <!-- same template, but now bound to Tailwind/shadcn classes via `ui.*` -->
  <Primitive :as="as" :class="ui.root({ class: [props.ui?.root, props.class] })">
    <table ref="tableRef" :class="ui.base({ class: [props.ui?.base] })">
      <caption v-if="caption || !!$slots.caption" :class="ui.caption({ class: [props.ui?.caption] })">
        <slot name="caption">
          {{ caption }}
        </slot>
      </caption>

      <thead :class="ui.thead({ class: [props.ui?.thead] })">
        <tr
          v-for="headerGroup in tableApi.getHeaderGroups()"
          :key="headerGroup.id"
          :class="ui.tr({ class: [props.ui?.tr] })"
        >
          <th
            v-for="header in headerGroup.headers"
            :key="header.id"
            :data-pinned="header.column.getIsPinned()"
            :scope="header.colSpan > 1 ? 'colgroup' : 'col'"
            :colspan="header.colSpan > 1 ? header.colSpan : undefined"
            :rowspan="header.rowSpan > 1 ? header.rowSpan : undefined"
            :class="
              ui.th({
                class: [props.ui?.th, resolveValue(header.column.columnDef.meta?.class?.th, header)],
                pinned: !!header.column.getIsPinned()
              })
            "
          >
            <slot :name="`${header.id}-header`" v-bind="header.getContext()">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </slot>
          </th>
        </tr>

        <tr :class="ui.separator({ class: [props.ui?.separator] })" />
      </thead>

      <tbody :class="ui.tbody({ class: [props.ui?.tbody] })">
        <slot name="body-top" />

        <template v-if="tableApi.getRowModel().rows?.length">
          <template v-for="row in tableApi.getRowModel().rows" :key="row.id">
            <tr
              :data-selected="row.getIsSelected()"
              :data-selectable="!!props.onSelect || !!props.onHover || !!props.onContextmenu"
              :data-expanded="row.getIsExpanded()"
              :role="props.onSelect ? 'button' : undefined"
              :tabindex="props.onSelect ? 0 : undefined"
              :class="
                ui.tr({
                  class: [props.ui?.tr, resolveValue(tableApi.options.meta?.class?.tr, row)]
                })
              "
              :style="resolveValue(tableApi.options.meta?.style?.tr, row)"
              @click="onRowSelect($event, row)"
              @pointerenter="onRowHover($event, row)"
              @pointerleave="onRowHover($event, null)"
              @contextmenu="onRowContextmenu($event, row)"
            >
              <td
                v-for="cell in row.getVisibleCells()"
                :key="cell.id"
                :data-pinned="cell.column.getIsPinned()"
                :colspan="resolveValue(cell.column.columnDef.meta?.colspan?.td, cell)"
                :rowspan="resolveValue(cell.column.columnDef.meta?.rowspan?.td, cell)"
                :class="
                  ui.td({
                    class: [props.ui?.td, resolveValue(cell.column.columnDef.meta?.class?.td, cell)],
                    pinned: !!cell.column.getIsPinned()
                  })
                "
                :style="resolveValue(cell.column.columnDef.meta?.style?.td, cell)"
              >
                <slot :name="`${cell.column.id}-cell`" v-bind="cell.getContext()">
                  <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                </slot>
              </td>
            </tr>
            <tr v-if="row.getIsExpanded()" :class="ui.tr({ class: [props.ui?.tr] })">
              <td :colspan="row.getAllCells().length" :class="ui.td({ class: [props.ui?.td] })">
                <slot name="expanded" :row="row" />
              </td>
            </tr>
          </template>
        </template>

        <tr v-else-if="loading && !!$slots['loading']">
          <td :colspan="tableApi.getAllLeafColumns().length" :class="ui.loading({ class: props.ui?.loading })">
            <slot name="loading" />
          </td>
        </tr>

        <tr v-else>
          <td :colspan="tableApi.getAllLeafColumns().length" :class="ui.empty({ class: props.ui?.empty })">
            <slot name="empty">
              {{ empty || "No Results" }}
            </slot>
          </td>
        </tr>

        <slot name="body-bottom" />
      </tbody>

      <tfoot v-if="hasFooter" :class="ui.tfoot({ class: [props.ui?.tfoot] })">
        <tr :class="ui.separator({ class: [props.ui?.separator] })" />
        <tr
          v-for="footerGroup in tableApi.getFooterGroups()"
          :key="footerGroup.id"
          :class="ui.tr({ class: [props.ui?.tr] })"
        >
          <th
            v-for="header in footerGroup.headers"
            :key="header.id"
            :data-pinned="header.column.getIsPinned()"
            :colspan="header.colSpan > 1 ? header.colSpan : undefined"
            :rowspan="header.rowSpan > 1 ? header.rowSpan : undefined"
            :class="
              ui.th({
                class: [props.ui?.th, resolveValue(header.column.columnDef.meta?.class?.th, header)],
                pinned: !!header.column.getIsPinned()
              })
            "
            :style="resolveValue(header.column.columnDef.meta?.style?.th, header)"
          >
            <slot :name="`${header.id}-footer`" v-bind="header.getContext()">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.footer"
                :props="header.getContext()"
              />
            </slot>
          </th>
        </tr>
      </tfoot>
    </table>
  </Primitive>
</template>
