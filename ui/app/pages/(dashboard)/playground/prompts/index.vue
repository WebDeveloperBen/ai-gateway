<template>
  <ClientOnly>
    <MonacoEditor
      v-model="val"
      lang="cel"
      :options="{ automaticLayout: true, minimap: { enabled: false }, wordWrap: 'on' }"
      @load="onLoad"
      ref="editorRef"
      style="height: 60vh"
    >
      Loading...
    </MonacoEditor>
  </ClientOnly>

  <div class="mt-3 flex gap-2">
    <button type="button" class="btn" @click="insertSnippet('noPiiProd')">No PII in prod</button>
    <button type="button" class="btn" @click="insertSnippet('auRegionsProd')">AU regions in prod</button>
    <button type="button" class="btn" @click="insertSnippet('nonProdCost')">Non-prod cost < 1¢</button>
  </div>
</template>

<script setup lang="ts">
import type { MonacoEditor } from "#components"
import type * as Monaco from "monaco-editor"

const val = ref<string | undefined>('cost < 0.01 && region == "australiaeast"')
const editorRef = useTemplateRef<InstanceType<typeof MonacoEditor>>("editorRef")

const SNIPPETS = {
  noPiiProd: '!(env == "prod" && pii_detected)',
  auRegionsProd: '!(env == "prod") || region.startsWith("australia")',
  nonProdCost: '(env != "prod") ? cost < 0.01 : true'
} as const

async function onLoad(editor: Monaco.editor.IStandaloneCodeEditor) {
  console.log("loaded", !!editor) // this should log
  // Optional: theme via useMonaco
  const monaco = await useMonaco()
  const colorMode = useColorMode()
  watch(
    () => colorMode.preference,
    (t) => monaco?.editor.setTheme(t === "dark" ? "vs-dark" : "vs"),
    { immediate: true }
  )
}

// bulletproof insertion using model edits
function insertSnippet(key: keyof typeof SNIPPETS) {
  const ed = editorRef.value?.$editor
  if (!ed) return
  const model = ed.getModel()
  if (!model) return

  const text = SNIPPETS[key]
  const sel = ed.getSelection()
  const pos = ed.getPosition()
  const range = sel ?? {
    startLineNumber: pos?.lineNumber ?? 1,
    startColumn: pos?.column ?? 1,
    endLineNumber: pos?.lineNumber ?? 1,
    endColumn: pos?.column ?? 1
  }

  model.pushEditOperations(sel ? [sel] : [], [{ range, text, forceMoveMarkers: true }], () => null)
  ed.focus()
}
</script>
