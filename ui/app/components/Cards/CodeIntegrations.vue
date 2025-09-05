<script setup lang="ts">
import { Code, Copy } from "lucide-vue-next"
import { toast } from "vue-sonner"

type CodeExampleKey = "curl" | "javascript" | "python"

const codeExamples: Record<CodeExampleKey, string> = {
  curl: `curl -X POST "https://api.yourdomain.com/v1/chat/completions" \\
  -H "Authorization: Bearer YOUR_API_KEY_HERE" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'`,

  javascript: `const response = await fetch('https://api.yourdomain.com/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY_HERE',
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    model: 'gpt-4',
    messages: [{ role: 'user', content: 'Hello!' }]
  })
});`,

  python: `import requests

headers = {
    'Authorization': 'Bearer YOUR_API_KEY_HERE',
    'Content-Type': 'application/json',
}

data = {
    'model': 'gpt-4',
    'messages': [{'role': 'user', 'content': 'Hello!'}]
}

response = requests.post('https://api.yourdomain.com/v1/chat/completions', 
                        headers=headers, json=data)`
}

const selectedCodeExample = ref<CodeExampleKey>("curl")

const copyCodeExample = async (example: CodeExampleKey) => {
  try {
    await navigator.clipboard.writeText(codeExamples[example])
    toast.success("Code copied successfully!", {})
    // Could add toast notification here
  } catch (err) {
    console.error("Failed to copy code example: ", err)
  }
}
</script>
<template>
  <UiCard>
    <UiCardHeader>
      <UiCardTitle class="flex items-center gap-2">
        <Code class="h-5 w-5" />
        Integration Examples
      </UiCardTitle>
      <UiCardDescription>Sample code to get you started with this API key</UiCardDescription>
    </UiCardHeader>
    <UiCardContent>
      <!-- Language Tabs -->
      <div class="flex items-center gap-2 mb-4">
        <UiButton
          v-for="lang in ['curl', 'javascript', 'python'] as CodeExampleKey[]"
          :key="lang"
          variant="outline"
          size="sm"
          @click="selectedCodeExample = lang"
          :class="selectedCodeExample === lang ? 'bg-primary text-primary-foreground' : ''"
        >
          {{ lang.charAt(0).toUpperCase() + lang.slice(1) }}
        </UiButton>
      </div>

      <!-- Code Block -->
      <div class="relative">
        <pre
          class="bg-muted p-4 rounded-lg text-sm overflow-x-auto"
        ><code>{{ codeExamples[selectedCodeExample]}}</code></pre>
        <UiButton
          variant="outline"
          size="sm"
          class="absolute top-2 right-2 gap-2"
          @click="copyCodeExample(selectedCodeExample)"
        >
          <Copy class="h-3 w-3" />
          Copy
        </UiButton>
      </div>
    </UiCardContent>
  </UiCard>
</template>
