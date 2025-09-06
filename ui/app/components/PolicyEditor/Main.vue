<template>
  <ClientOnly>
    <div class="flex flex-col gap-6 h-full">
      <!-- Policy Configuration Form -->
      <form @submit.prevent>
        <fieldset :disabled="isSubmitting">
          <UiFormBuilder class="grid grid-cols-12 gap-5" :fields="formFields" />
        </fieldset>
      </form>

      <!-- Main Editor Area -->
      <div class="flex gap-6 flex-1 min-h-0">
        <!-- CEL Editor -->
        <div class="flex-1 flex flex-col">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <Code2 class="h-4 w-4 text-primary" />
              <h3 class="font-medium text-sm">CEL Expression</h3>
            </div>

            <!-- Status Indicator and Controls -->
            <div class="flex items-center gap-2">
              <button
                class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-md bg-secondary text-secondary-foreground hover:bg-secondary/80 border border-border/50 transition-all duration-150 ease-in-out hover:shadow-sm hover:border-border disabled:opacity-50 disabled:cursor-not-allowed"
                type="button"
                @click="validate"
                :disabled="validating"
                title="Validate CEL expression"
              >
                <CheckCircle2 v-if="!validating" class="h-3.5 w-3.5" />
                <Loader2 v-else class="h-3.5 w-3.5 animate-spin" />
                <span class="text-xs">Validate</span>
              </button>

              <button
                class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-md border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out hover:shadow-sm hover:border-border/60 disabled:opacity-50 disabled:cursor-not-allowed"
                type="button"
                @click="draftWithAI"
                :disabled="drafting"
                title="Generate policy with AI"
              >
                <Sparkles v-if="!drafting" class="h-3.5 w-3.5" />
                <Loader2 v-else class="h-3.5 w-3.5 animate-spin" />
                <span class="text-xs">AI Draft</span>
              </button>

              <div class="flex items-center gap-2 ml-2">
                <div
                  class="h-2 w-2 rounded-full transition-colors"
                  :class="{
                    'bg-green-500': validationStatus === 'valid',
                    'bg-red-500': validationStatus === 'invalid',
                    'bg-yellow-500': validationStatus === 'validating',
                    'bg-gray-400': validationStatus === 'unknown'
                  }"
                />
                <span class="text-xs text-muted-foreground">{{ validationStatusText }}</span>
                <button
                  v-if="proposed"
                  class="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out"
                  type="button"
                  @click="openDiff"
                  title="Show proposed changes"
                >
                  <Diff class="h-3 w-3" />
                  Diff
                </button>
              </div>
            </div>
          </div>

          <div class="flex-1 border border-border rounded-lg">
            <MonacoEditor
              ref="editorRef"
              v-model="value"
              lang="cel"
              :options="editorOptions"
              @load="onLoad"
              class="border-0"
              style="height: 100%"
            >
              <div class="flex h-full items-center justify-center">
                <div class="flex flex-col items-center gap-2 text-muted-foreground">
                  <Loader2 class="h-6 w-6 animate-spin" />
                  <p class="text-sm">Loading Monaco Editor...</p>
                </div>
              </div>
            </MonacoEditor>

            <!-- Editor overlay for additional info -->
            <div
              v-if="hasErrors"
              class="absolute bottom-3 right-3 flex items-center gap-2 rounded-md bg-red-50 dark:bg-red-950/20 px-2 py-1 text-xs z-20"
            >
              <AlertCircle class="h-3 w-3 text-red-500" />
              <span class="text-red-700 dark:text-red-400">{{ errorCount }} error{{ errorCount > 1 ? "s" : "" }}</span>
            </div>
          </div>
        </div>

        <!-- Template Policies Sidebar -->
        <div class="w-80 rounded-lg border bg-background shadow-sm">
          <div class="p-4 border-b">
            <h3 class="font-semibold text-sm flex items-center gap-2">
              <Library class="h-4 w-4 text-primary" />
              Policy Templates
            </h3>
            <p class="text-xs text-muted-foreground mt-1">Click to insert into editor</p>
          </div>

          <div class="p-3 max-h-[60vh] overflow-y-auto">
            <div class="grid grid-cols-2 gap-3">
              <div
                v-for="template in policyTemplates"
                :key="template.key"
                class="flex flex-col gap-2 p-3 rounded-lg border border-border bg-card hover:bg-muted/30 cursor-pointer transition-all duration-150 ease-in-out group hover:shadow-sm hover:border-border/60"
                @click="insertSnippet(template.key)"
                :title="template.description"
              >
                <div
                  class="flex items-center justify-center h-8 w-8 rounded-md bg-primary/10 group-hover:bg-primary/20 transition-colors"
                >
                  <component :is="template.icon" class="h-4 w-4 text-primary" />
                </div>
                <div>
                  <p class="text-xs font-medium text-foreground">{{ template.name }}</p>
                  <p class="text-xs text-muted-foreground leading-tight">{{ template.shortDesc }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Bar -->
      <div class="flex items-center justify-end pt-4 gap-3">
        <button
          class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out hover:shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          type="button"
          @click="saveDraft"
          :disabled="isSubmitting"
        >
          <Save class="h-4 w-4" />
          Save Draft
        </button>

        <button
          class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-all duration-150 ease-in-out hover:shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          type="button"
          @click="validateAndDeploy"
          :disabled="isSubmitting || !value || validationStatus === 'invalid'"
        >
          <CheckCircle2 class="h-4 w-4" />
          Validate & Deploy
        </button>
      </div>
    </div>

    <!-- Diff modal -->
    <div v-if="showDiff" class="fixed inset-0 z-50 grid place-items-center bg-black/40 p-4">
      <div class="w-full max-w-5xl rounded-xl border bg-card shadow-xl">
        <div class="flex items-center justify-between border-b px-3 py-2">
          <div class="text-sm font-medium">Proposed change vs Current</div>
          <div class="flex items-center gap-2">
            <button
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-all duration-150 ease-in-out hover:shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
              type="button"
              @click="acceptProposed"
              :disabled="!proposed"
            >
              Accept
            </button>
            <button
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium rounded-md border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out hover:shadow-sm hover:border-border/60"
              type="button"
              @click="closeDiff"
            >
              Close
            </button>
          </div>
        </div>
        <MonacoDiffEditor
          ref="diffRef"
          :original="value ?? ''"
          :value="proposed ?? ''"
          :lang="'cel'"
          :options="{ renderSideBySide: true, automaticLayout: true, minimap: { enabled: false } }"
          @load="onDiffLoad"
          style="height: 60vh"
        />
      </div>
    </div>
  </ClientOnly>
</template>

<script setup lang="ts">
import type { MonacoEditor, MonacoDiffEditor } from "#components"
import type * as MonacoNS from "monaco-editor"
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue"
import { toast } from "vue-sonner"
import * as z from "zod"
import {
  Code2,
  Shield,
  Globe,
  DollarSign,
  CheckCircle2,
  Sparkles,
  Diff,
  Loader2,
  AlertCircle,
  Library,
  Clock,
  Users,
  Database,
  Zap,
  Lock,
  Save
} from "lucide-vue-next"

// Form Schema
const policySchema = z.object({
  name: z.string().min(1, "Policy name is required"),
  description: z.string().optional(),
  targetEnvironment: z
    .array(z.enum(["all", "dev", "staging", "prod"]))
    .min(1, "At least one environment must be selected"),
  applications: z.array(z.string()).min(1, "At least one application must be selected")
})

// Form setup
const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(policySchema),
  initialValues: {
    name: "Untitled Policy",
    description: "",
    targetEnvironment: ["all"] as const,
    applications: [] as string[]
  }
})

// Mock data - replace with real data from your API
const availableApps = ref([
  { label: "Chat Application", value: "chat-app" },
  { label: "Analytics API", value: "analytics-api" },
  { label: "User Service", value: "user-service" },
  { label: "Content API", value: "content-api" },
  { label: "Recommendation Engine", value: "recommendation-engine" }
])

/** -------------------------------
 *  1) Your theme HEX palette (fill these)
 *  -------------------------------- */
const LIGHT = {
  bg: "#f9f9f9",
  fg: "#1f2937",
  gutter: "#f9f9f9",
  line: "#6b7280",
  sel: "#e5e7eb",
  cursor: "#111827"
}
const DARK = {
  bg: "#18181a",
  fg: "#e5e7eb",
  gutter: "#0f162a",
  line: "#94a3b8",
  sel: "#334155",
  cursor: "#e5e7eb"
}

// Form fields definition
const formFields: FormBuilder[] = [
  {
    variant: "Input",
    label: "Policy Name",
    name: "name",
    placeholder: "e.g. PII Protection Policy",
    required: true,
    wrapperClass: "col-span-full md:col-span-6 lg:col-span-3"
  },
  {
    variant: "MultiSelect",
    label: "Environment",
    name: "targetEnvironment",
    placeholder: "Select environments...",
    required: true,
    mode: "multiple",
    hideSelected: false,
    searchable: true,
    closeOnSelect: false,
    closeOnDeselect: false,
    canClear: true,
    options: [
      { label: "All Environments", value: "all" },
      { label: "Development", value: "dev" },
      { label: "Staging", value: "staging" },
      { label: "Production", value: "prod" }
    ],
    classes: {
      option:
        "flex items-center gap-2 px-3 py-2 hover:bg-accent hover:text-accent-foreground cursor-pointer transition-colors",
      optionSelected: "bg-accent/50 text-accent-foreground",
      optionPointed: "bg-accent text-accent-foreground"
    },
    wrapperClass: "col-span-full md:col-span-6 lg:col-span-3"
  },
  {
    variant: "MultiSelect",
    label: "Applications",
    name: "applications",
    placeholder: "Select applications...",
    required: true,
    mode: "multiple",
    hideSelected: false,
    searchable: true,
    closeOnSelect: false,
    closeOnDeselect: false,
    canClear: true,
    options: availableApps.value,
    classes: {
      option:
        "flex items-center gap-2 px-3 py-2 hover:bg-accent hover:text-accent-foreground cursor-pointer transition-colors",
      optionSelected: "bg-accent/50 text-accent-foreground",
      optionPointed: "bg-accent text-accent-foreground"
    },
    wrapperClass: "col-span-full md:col-span-6 lg:col-span-3"
  },
  {
    variant: "Textarea",
    label: "Policy Description",
    name: "description",
    placeholder: "Describe what this policy does and when it applies...",
    rows: 2,
    wrapperClass: "col-span-full"
  }
]

/** -------------------------------
 *  2) State & refs
 *  -------------------------------- */
const value = ref<string | undefined>('cost < 0.01 && region == "australiaeast"')
const proposed = ref<string | undefined>(undefined) // AI suggestion / draft
const showDiff = ref(false)

// UI State
const validating = ref(false)
const drafting = ref(false)
const validationStatus = ref<"valid" | "invalid" | "validating" | "unknown">("unknown")
const currentErrors = ref<
  Array<{ message: string; from: { line: number; col: number }; to: { line: number; col: number } }>
>([])

const editorRef = useTemplateRef<InstanceType<typeof MonacoEditor>>("editorRef")
const diffRef = useTemplateRef<InstanceType<typeof MonacoDiffEditor>>("diffRef")

let editor: MonacoNS.editor.IStandaloneCodeEditor | null = null
let monaco: typeof import("monaco-editor") | null = null

// Computed properties
const hasErrors = computed(() => currentErrors.value.length > 0)
const errorCount = computed(() => currentErrors.value.length)
const validationStatusText = computed(() => {
  switch (validationStatus.value) {
    case "valid":
      return "Valid"
    case "invalid":
      return "Invalid"
    case "validating":
      return "Validating..."
    default:
      return "Ready"
  }
})

// Form submission handlers
const saveDraft = handleSubmit(async (formData) => {
  try {
    const policyData = {
      ...formData,
      celExpression: value.value,
      status: "draft"
    }

    console.log("Saving draft policy:", policyData)
    // await $fetch('/api/policies/draft', { method: 'POST', body: policyData })

    toast.success("Draft saved successfully!")
  } catch (error) {
    console.error("Failed to save draft:", error)
    toast.error("Failed to save draft")
  }
})

const validateAndDeploy = handleSubmit(async (formData) => {
  try {
    // First validate CEL
    await validate()

    if (validationStatus.value !== "valid") {
      toast.error("Please fix validation errors before deploying")
      return
    }

    const policyData = {
      ...formData,
      celExpression: value.value,
      status: "active"
    }

    console.log("Deploying policy:", policyData)
    // await $fetch('/api/policies/deploy', { method: 'POST', body: policyData })

    toast.success("Policy deployed successfully!")
  } catch (error) {
    console.error("Failed to deploy policy:", error)
    toast.error("Failed to deploy policy")
  }
})

/** -------------------------------
 *  3) Editor options / snippets
 *  -------------------------------- */
const editorOptions: MonacoNS.editor.IStandaloneEditorConstructionOptions = {
  automaticLayout: true,
  minimap: { enabled: false },
  wordWrap: "on",
  fontLigatures: true,
  fontSize: 14,
  lineHeight: 20,
  padding: { top: 0, bottom: 16 },
  scrollBeyondLastLine: false,
  renderLineHighlight: "line",
  smoothScrolling: true,
  cursorBlinking: "smooth",
  acceptSuggestionOnEnter: "smart",
  "semanticHighlighting.enabled": true,
  codeLens: true,
  colorDecorators: true,
  inlineSuggest: {
    enabled: true,
    keepOnBlur: false
  },
  showUnused: true,
  bracketPairColorization: { enabled: true },
  hover: {
    enabled: true,
    delay: 100,
    sticky: true
  },
  quickSuggestions: {
    other: true,
    comments: true,
    strings: true
  },
  parameterHints: {
    enabled: true
  }
}

const SNIPPETS = {
  noPiiProd: '!(env == "prod" && pii_detected)',
  auRegionsProd: '!(env == "prod") || region.startsWith("australia")',
  nonProdCost: '(env != "prod") ? cost < 0.01 : true',
  rateLimiting: "tokens <= 10000 && cost <= 1.0",
  officeHours: "hour >= 9 && hour <= 17",
  userRole: 'role == "admin" || role == "developer"',
  modelAccess: 'model.startsWith("gpt-") && cost <= 5.0',
  dataClassification: "!pii_detected || redaction_applied"
} as const

// Enhanced policy templates for sidebar
const policyTemplates = [
  {
    key: "noPiiProd",
    name: "No PII in Prod",
    shortDesc: "Block PII data in production",
    description: "Prevents PII-containing requests in production environments",
    icon: Shield
  },
  {
    key: "auRegionsProd",
    name: "AU Regions",
    shortDesc: "Restrict to AU regions",
    description: "Allows only Australian regions for production traffic",
    icon: Globe
  },
  {
    key: "nonProdCost",
    name: "Cost Control",
    shortDesc: "Limit non-prod costs",
    description: "Controls costs for non-production environments",
    icon: DollarSign
  },
  {
    key: "rateLimiting",
    name: "Rate Limits",
    shortDesc: "Token & cost limits",
    description: "Enforces rate limiting based on tokens and cost",
    icon: Zap
  },
  {
    key: "officeHours",
    name: "Office Hours",
    shortDesc: "Business hours only",
    description: "Restricts access to business hours (9-5)",
    icon: Clock
  },
  {
    key: "userRole",
    name: "Role Access",
    shortDesc: "Admin/dev only",
    description: "Restricts access to admin and developer roles",
    icon: Users
  },
  {
    key: "modelAccess",
    name: "Model Control",
    shortDesc: "GPT models + cost",
    description: "Allows GPT models with cost constraints",
    icon: Database
  },
  {
    key: "dataClassification",
    name: "Data Safety",
    shortDesc: "PII redaction",
    description: "Ensures PII is redacted when detected",
    icon: Lock
  }
] as const

function insertSnippet(key: keyof typeof SNIPPETS) {
  const ed = editorRef.value?.$editor
  if (!ed) return
  const model = ed.getModel()
  if (!model) return
  const sel = ed.getSelection()
  const pos = ed.getPosition()
  const range = sel ?? {
    startLineNumber: pos?.lineNumber ?? 1,
    startColumn: pos?.column ?? 1,
    endLineNumber: pos?.lineNumber ?? 1,
    endColumn: pos?.column ?? 1
  }
  model.pushEditOperations(sel ? [sel] : [], [{ range, text: SNIPPETS[key], forceMoveMarkers: true }], () => null)
  ed.focus()
}

/** -------------------------------
 *  4) Custom CEL Language Definition
 *  -------------------------------- */

// Define allowed variables and functions for your CEL policies
const CEL_VARIABLES = [
  "role",
  "app",
  "model",
  "region",
  "env",
  "tokens",
  "cost",
  "pii_detected",
  "redaction_applied",
  "hour",
  "user_id",
  "request_size",
  "response_size",
  "latency"
]

const CEL_FUNCTIONS = [
  "startsWith",
  "endsWith",
  "contains",
  "matches",
  "size",
  "has",
  "in",
  "all",
  "exists",
  "exists_one",
  "map",
  "filter",
  "toString",
  "int",
  "double"
]

const CEL_OPERATORS = [
  "==",
  "!=",
  "<",
  "<=",
  ">",
  ">=",
  "&&",
  "||",
  "!",
  "+",
  "-",
  "*",
  "/",
  "%",
  "in",
  "not",
  "and",
  "or"
]

function registerCELLanguage(m: typeof import("monaco-editor")) {
  // Don't re-register if already exists
  const existingLanguages = m.languages.getLanguages()
  if (!existingLanguages.find((lang) => lang.id === "cel")) {
    m.languages.register({ id: "cel" })
  }

  // Add autocompletion for CEL
  m.languages.registerCompletionItemProvider("cel", {
    provideCompletionItems: (model, position) => {
      const word = model.getWordUntilPosition(position)
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn
      }

      const suggestions = [
        // Variables
        ...CEL_VARIABLES.map((variable) => ({
          label: variable,
          kind: m.languages.CompletionItemKind.Variable,
          insertText: variable,
          range,
          documentation: getVariableDocumentation(variable),
          detail: `Variable: ${variable}`
        })),
        // Functions
        ...CEL_FUNCTIONS.map((func) => ({
          label: func,
          kind: m.languages.CompletionItemKind.Function,
          insertText: `${func}()`,
          range,
          documentation: getFunctionDocumentation(func),
          detail: `Function: ${func}()`
        })),
        // Constants
        { label: "true", kind: m.languages.CompletionItemKind.Constant, insertText: "true", range },
        { label: "false", kind: m.languages.CompletionItemKind.Constant, insertText: "false", range },
        { label: "null", kind: m.languages.CompletionItemKind.Constant, insertText: "null", range },
        // Common operators
        { label: "==", kind: m.languages.CompletionItemKind.Operator, insertText: "== ", range },
        { label: "!=", kind: m.languages.CompletionItemKind.Operator, insertText: "!= ", range },
        { label: "&&", kind: m.languages.CompletionItemKind.Operator, insertText: "&& ", range },
        { label: "||", kind: m.languages.CompletionItemKind.Operator, insertText: "|| ", range }
      ]

      return { suggestions }
    }
  })

  // Add hover documentation - restore working version
  m.languages.registerHoverProvider("cel", {
    provideHover: (model, position) => {
      // First check if there are any markers at this position
      const markers = m.editor.getModelMarkers({ resource: model.uri })
      const markersAtPosition = markers.filter(
        (marker) =>
          marker.startLineNumber === position.lineNumber &&
          position.column >= marker.startColumn &&
          position.column <= marker.endColumn
      )

      if (markersAtPosition.length > 0) {
        const firstMarker = markersAtPosition[0]!
        return {
          range: new m.Range(
            firstMarker.startLineNumber,
            firstMarker.startColumn,
            firstMarker.endLineNumber,
            firstMarker.endColumn
          ),
          contents: markersAtPosition.map((marker) => ({ value: `❌ ${marker.message}` }))
        }
      }

      // Then check for variable/function documentation
      const word = model.getWordAtPosition(position)
      if (!word) return null

      const { word: wordText } = word

      if (CEL_VARIABLES.includes(wordText)) {
        return {
          range: new m.Range(position.lineNumber, word.startColumn, position.lineNumber, word.endColumn),
          contents: [{ value: `**${wordText}** (variable)` }, { value: getVariableDocumentation(wordText) }]
        }
      }

      if (CEL_FUNCTIONS.includes(wordText)) {
        return {
          range: new m.Range(position.lineNumber, word.startColumn, position.lineNumber, word.endColumn),
          contents: [{ value: `**${wordText}** (function)` }, { value: getFunctionDocumentation(wordText) }]
        }
      }

      return null
    }
  })

  // Add quick-fix code actions
  m.languages.registerCodeActionProvider("cel", {
    provideCodeActions: (model, _range, context) => {
      const actions: any[] = []

      for (const marker of context.markers) {
        if (marker.source !== "CEL Lint") continue

        // Fix single = to ==
        if (marker.message.includes("Use '==' for comparison")) {
          actions.push({
            title: "Replace '=' with '=='",
            kind: "quickfix",
            edit: {
              edits: [
                {
                  resource: model.uri,
                  textEdit: {
                    range: new m.Range(
                      marker.startLineNumber,
                      marker.startColumn,
                      marker.endLineNumber,
                      marker.endColumn
                    ),
                    text: "=="
                  }
                }
              ]
            },
            isPreferred: true
          })
        }

        // Fix typos in variable names
        if (marker.message.includes("Did you mean")) {
          const suggestion = marker.message.match(/'([^']+)'/g)?.[1]?.replace(/'/g, "")
          if (suggestion) {
            actions.push({
              title: `Replace with '${suggestion}'`,
              kind: "quickfix",
              edit: {
                edits: [
                  {
                    resource: model.uri,
                    textEdit: {
                      range: new m.Range(
                        marker.startLineNumber,
                        marker.startColumn,
                        marker.endLineNumber,
                        marker.endColumn
                      ),
                      text: suggestion
                    }
                  }
                ]
              },
              isPreferred: true
            })
          }
        }
      }

      return {
        actions,
        dispose: () => {}
      }
    }
  })
}

function getVariableDocumentation(variable: string): string {
  const docs: Record<string, string> = {
    role: 'User role (e.g., "admin", "developer", "user")',
    app: "Application identifier",
    model: 'AI model being used (e.g., "gpt-4", "claude-3")',
    region: 'Geographic region (e.g., "us-east-1", "australiaeast")',
    env: 'Environment (e.g., "prod", "staging", "dev")',
    tokens: "Number of tokens in the request",
    cost: "Estimated cost of the request",
    pii_detected: "Boolean indicating if PII was detected",
    redaction_applied: "Boolean indicating if redaction was applied",
    hour: "Current hour (0-23)",
    user_id: "Unique user identifier",
    request_size: "Size of the request in bytes",
    response_size: "Size of the response in bytes",
    latency: "Request latency in milliseconds"
  }
  return docs[variable] || "Custom variable"
}

function getFunctionDocumentation(func: string): string {
  const docs: Record<string, string> = {
    startsWith: "startsWith(string, prefix) - Check if string starts with prefix",
    endsWith: "endsWith(string, suffix) - Check if string ends with suffix",
    contains: "contains(string, substring) - Check if string contains substring",
    matches: "matches(string, regex) - Check if string matches regex pattern",
    size: "size(collection) - Get size of string, list, or map",
    has: "has(object, field) - Check if object has field",
    in: "value in collection - Check if value exists in collection",
    all: "all(list, predicate) - Check if all items match predicate",
    exists: "exists(list, predicate) - Check if any item matches predicate"
  }
  return docs[func] || "Built-in function"
}

/** -------------------------------
 *  5) Theme registration (light/dark)
 *  -------------------------------- */
function defineThemes(m: typeof import("monaco-editor")) {
  // Enhanced CEL syntax highlighting
  m.editor.defineTheme("cel-light", {
    base: "vs",
    inherit: true,
    rules: [
      { token: "keyword", foreground: "0969da", fontStyle: "bold" },
      { token: "identifier", foreground: "24292f" },
      { token: "number", foreground: "0550ae" },
      { token: "string", foreground: "0a3069" },
      { token: "operator", foreground: "cf222e", fontStyle: "bold" },
      { token: "delimiter", foreground: "6f42c1" },
      { token: "comment", foreground: "656d76", fontStyle: "italic" },
      { token: "function", foreground: "8250df", fontStyle: "bold" },
      { token: "variable", foreground: "953800" }
    ],
    colors: {
      "editor.background": LIGHT.bg,
      "editor.foreground": LIGHT.fg,
      "editorGutter.background": LIGHT.gutter,
      "editorLineNumber.foreground": LIGHT.line,
      "editorLineNumber.activeForeground": LIGHT.fg,
      "editor.selectionBackground": LIGHT.sel,
      "editor.inactiveSelectionBackground": LIGHT.sel,
      "editorCursor.foreground": LIGHT.cursor,
      "scrollbarSlider.background": LIGHT.sel,
      "editor.lineHighlightBackground": "#f6f8fa"
    }
  })

  m.editor.defineTheme("cel-dark", {
    base: "vs-dark",
    inherit: true,
    rules: [
      { token: "keyword", foreground: "79c0ff", fontStyle: "bold" },
      { token: "identifier", foreground: "e6edf3" },
      { token: "number", foreground: "a5d6ff" },
      { token: "string", foreground: "a5d6ff" },
      { token: "operator", foreground: "ff7b72", fontStyle: "bold" },
      { token: "delimiter", foreground: "d2a8ff" },
      { token: "comment", foreground: "8b949e", fontStyle: "italic" },
      { token: "function", foreground: "d2a8ff", fontStyle: "bold" },
      { token: "variable", foreground: "ffa657" }
    ],
    colors: {
      "editor.background": DARK.bg,
      "editor.foreground": DARK.fg,
      "editorGutter.background": DARK.gutter,
      "editorLineNumber.foreground": DARK.line,
      "editorLineNumber.activeForeground": DARK.fg,
      "editor.selectionBackground": DARK.sel,
      "editor.inactiveSelectionBackground": DARK.sel,
      "editorCursor.foreground": DARK.cursor,
      "scrollbarSlider.background": DARK.sel,
      "editor.lineHighlightBackground": "#21262d"
    }
  })
}

function applyThemeForColorMode() {
  const cm = useColorMode()
  const theme = cm.preference === "dark" ? "cel-dark" : "cel-light"
  monaco?.editor.setTheme(theme)
}

/** -------------------------------
 *  5) Validation (server truth)
 *  -------------------------------- */
type ValidateRes = {
  ok: boolean
  errors?: Array<{ message: string; from: { line: number; col: number }; to: { line: number; col: number } }>
}

async function validate() {
  validating.value = true
  validationStatus.value = "validating"

  try {
    const res = await $fetch<ValidateRes>("/api/cel/validate", {
      method: "POST",
      body: { expr: value.value ?? "" }
    })
    const errors = res.errors ?? []
    currentErrors.value = errors
    setMarkers(errors)

    validationStatus.value = errors.length > 0 ? "invalid" : "valid"
  } catch (error) {
    console.error("Validation failed:", error)
    currentErrors.value = []
    setMarkers([])
    validationStatus.value = "unknown"
  } finally {
    validating.value = false
  }
}

function setMarkers(
  errors: Array<{ message: string; from: { line: number; col: number }; to: { line: number; col: number } }>
) {
  const ed = editorRef.value?.$editor
  if (!ed || !monaco) return
  const model = ed.getModel()
  if (!model) return
  monaco.editor.setModelMarkers(
    model,
    "cel-validate",
    errors.map((e) => ({
      message: e.message,
      severity: monaco!.MarkerSeverity.Error,
      startLineNumber: e.from.line,
      startColumn: e.from.col,
      endLineNumber: e.to.line,
      endColumn: e.to.col
    }))
  )
}
const debouncedValidate = useDebounce(validate, 450)

/** -------------------------------
 *  6) AI draft stub (proposed text)
 *  -------------------------------- */
async function draftWithAI() {
  drafting.value = true

  try {
    const res = await $fetch<{ cel: string }>("/api/policy/ai-draft", {
      method: "POST",
      body: {
        instruction: "Block PII in prod, allow elsewhere",
        variables: ["role", "app", "model", "region", "env", "tokens", "cost", "pii_detected", "redaction_applied"]
      }
    })
    proposed.value = res.cel
    showDiff.value = true
  } catch (error) {
    console.error("AI draft failed:", error)
  } finally {
    drafting.value = false
  }
}

/** -------------------------------
 *  7) Diff controls
 *  -------------------------------- */
function openDiff() {
  showDiff.value = true
}
function closeDiff() {
  showDiff.value = false
}
function acceptProposed() {
  if (proposed.value) value.value = proposed.value
  showDiff.value = false
}

/** -------------------------------
 *  8) Lifecycle / load hooks
 *  -------------------------------- */
async function onLoad(ed: MonacoNS.editor.IStandaloneCodeEditor) {
  editor = ed
  monaco = await useMonaco()
  if (!monaco) return

  // Register custom CEL language
  registerCELLanguage(monaco)

  // themes
  defineThemes(monaco)
  applyThemeForColorMode()

  // react to light/dark toggles
  const cm = useColorMode()
  watch(
    () => cm.preference,
    () => applyThemeForColorMode()
  )

  // Add smart CEL linting
  const performCELLinting = () => {
    if (!editor || !monaco) return

    const model = editor.getModel()
    if (!model) return

    // Clear all existing markers first
    monaco.editor.setModelMarkers(model, "cel-linter", [])

    const text = model.getValue()
    console.log("Validating:", text)

    const markers: any[] = []

    // Check for single = instead of ==
    const lines = text.split("\n")
    lines.forEach((line, lineIndex) => {
      const lineNumber = lineIndex + 1

      // Look for single = (not ==)
      const singleEqualsMatch = line.match(/(\w+)\s*(=)\s*(?!=)/g)
      if (singleEqualsMatch) {
        const equalIndex = line.indexOf("=")
        if (equalIndex !== -1 && line[equalIndex + 1] !== "=") {
          markers.push({
            message: "Use '==' for comparison, not '='. CEL uses '==' for equality checks.",
            severity: 8, // Error
            startLineNumber: lineNumber,
            startColumn: equalIndex + 1,
            endLineNumber: lineNumber,
            endColumn: equalIndex + 2
          })
        }
      }

      // Check for invalid operators
      const operatorRegex = /([<>!=]+|[&|]{1,2}|[+\-*/%])/g
      let operatorMatch
      while ((operatorMatch = operatorRegex.exec(line)) !== null) {
        const operator = operatorMatch[1]
        if (!operator) return
        if (!CEL_OPERATORS.includes(operator)) {
          markers.push({
            message: `Invalid operator '${operator}'. Valid operators: ${CEL_OPERATORS.join(", ")}`,
            severity: 8, // Error
            startLineNumber: lineNumber,
            startColumn: operatorMatch.index + 1,
            endLineNumber: lineNumber,
            endColumn: operatorMatch.index + operator.length + 1
          })
        }
      }

      // Check for unknown variables
      const words = line.match(/\b[a-zA-Z_]\w*\b/g) || []
      words.forEach((word) => {
        // Skip if it's a known variable, function, keyword, or constant
        if (
          CEL_VARIABLES.includes(word) ||
          CEL_FUNCTIONS.includes(word) ||
          ["true", "false", "null", "if", "else", "in", "not", "and", "or"].includes(word)
        ) {
          return
        }

        // Check if it's in quotes (string literal) - improved logic
        const wordIndex = line.indexOf(word)
        const beforeWord = line.substring(0, wordIndex)

        // Count quotes before the word
        const quotesBefore = (beforeWord.match(/"/g) || []).length
        const singleQuotesBefore = (beforeWord.match(/'/g) || []).length

        // Check if we're inside quotes
        const inDoubleQuotes = quotesBefore % 2 === 1
        const inSingleQuotes = singleQuotesBefore % 2 === 1
        const inQuotes = inDoubleQuotes || inSingleQuotes

        if (!inQuotes) {
          markers.push({
            message: `Unknown variable '${word}'. Available: ${CEL_VARIABLES.slice(0, 5).join(", ")}...`,
            severity: 8, // Error
            startLineNumber: lineNumber,
            startColumn: wordIndex + 1,
            endLineNumber: lineNumber,
            endColumn: wordIndex + word.length + 1
          })
        }
      })
    })

    console.log("Setting markers:", markers.length, "errors found")

    // Set the new markers (we already cleared them at the start)
    monaco.editor.setModelMarkers(model, "cel-linter", markers)
  }

  // validate as you type (both server and client-side)
  editor.onDidChangeModelContent(() => {
    // Run validation immediately without debouncing for now
    performCELLinting()
    debouncedValidate()
  })

  // Add mouse hover listener as backup

  // Initial validation
  validate()
  performCELLinting()
}

function onDiffLoad(_ed: MonacoNS.editor.IStandaloneDiffEditor) {
  // Optional: diff-specific settings can go here
}

/** -------------------------------
 *  9) Small utils
 *  -------------------------------- */
function useDebounce<T extends (...a: any[]) => any>(fn: T, wait = 400) {
  let t: number | undefined
  return (...args: Parameters<T>) => {
    if (t) clearTimeout(t)
    // @ts-expect-error browser
    t = setTimeout(() => fn(...args), wait)
  }
}
</script>
