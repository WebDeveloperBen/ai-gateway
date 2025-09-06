<template>
   
  <MonacoEditor v-model="val" :lang="EXAMPLE_LANG_ID" :options="options" style="height: 100vh">
        Loading...  
  </MonacoEditor>
</template>

<script lang="ts" setup>
import { MonacoEditor } from "#components"
import type * as Monaco from "monaco-editor"
const colorMode = useColorMode()
if (import.meta.client) {
  const monaco = await useMonaco()
  watch(
    () => colorMode.preference,
    (theme) => monaco.editor.setTheme(theme === "dark" ? "vs-dark" : "vs"),
    { immediate: true }
  )
}
const EXAMPLE_LANG_ID = "examplelang"
const options: Ref<Monaco.editor.IStandaloneEditorConstructionOptions> = ref({
  automaticLayout: true
})

const val = ref("")

if (import.meta.client) {
  const monaco = await useMonaco()
  monaco.languages.register({
    id: EXAMPLE_LANG_ID
  })
  monaco.languages.registerCompletionItemProvider(EXAMPLE_LANG_ID, {
    triggerCharacters: ["."],
    provideCompletionItems: async function (_model, position, _context) {
      return {
        incomplete: true,
        suggestions: [
          {
            label: "someProp",
            kind: monaco.languages.CompletionItemKind.Property,
            insertText: "someProp",
            range: {
              startLineNumber: position.lineNumber,
              startColumn: position.column,
              endLineNumber: position.lineNumber,
              endColumn: position.column
            }
          }
        ]
      }
    }
  })
}
</script>
