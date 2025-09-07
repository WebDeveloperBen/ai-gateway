<script setup lang="ts">
import { Library, Settings, MessageSquare } from "lucide-vue-next"

interface Props {
  systemPromptTemplates: PromptTemplate[]
  userPromptTemplates: PromptTemplate[]
}

defineProps<Props>()

const playgroundState = usePlaygroundState()

if (!playgroundState) {
  throw new Error("Playground sidebar must be used within a playground state provider")
}

const { activeTemplateTab, activePromptTab, insertTemplate } = playgroundState
</script>

<template>
  <div class="w-80 flex flex-col rounded-xl border bg-card shadow-lg overflow-hidden">
    <!-- Templates Header -->
    <div class="p-4 border-b bg-card">
      <h3 class="font-semibold text-sm flex items-center gap-2 text-foreground">
        <Library class="h-4 w-4 text-primary" />
        Prompt Templates
      </h3>
      <p class="text-xs text-muted-foreground mt-1">Click to insert at cursor position</p>
    </div>

    <!-- Template Type Tabs -->
    <div class="border-b bg-muted/20">
      <div class="flex">
        <button
          @click="
            () => {
              activeTemplateTab = 'system'
              activePromptTab = 'system'
            }
          "
          :class="[
            'flex-1 px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center justify-center gap-2',
            activeTemplateTab === 'system'
              ? 'border-primary text-primary bg-background'
              : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
          ]"
        >
          <Settings class="h-4 w-4" />
          System
        </button>
        <button
          @click="
            () => {
              activeTemplateTab = 'user'
              activePromptTab = 'user'
            }
          "
          :class="[
            'flex-1 px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center justify-center gap-2',
            activeTemplateTab === 'user'
              ? 'border-primary text-primary bg-background'
              : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
          ]"
        >
          <MessageSquare class="h-4 w-4" />
          User
        </button>
      </div>
    </div>

    <!-- Template Content -->
    <div class="flex-1 min-h-0">
      <UiScrollArea class="h-full">
        <div class="p-4">
          <!-- System Templates -->
          <div v-show="activeTemplateTab === 'system'" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div
              v-for="template in systemPromptTemplates"
              :key="template.id"
              class="flex flex-col gap-3 p-4 rounded-xl border border-border bg-background hover:bg-primary/5 hover:border-primary/30 cursor-pointer transition-all duration-200 ease-in-out group hover:shadow-md hover:scale-[1.02]"
              @click="insertTemplate(template, 'system')"
              :title="template.description"
            >
              <div
                class="flex items-center justify-center h-10 w-10 rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors shadow-sm"
              >
                <component :is="template.icon" class="h-5 w-5 text-primary" />
              </div>
              <div>
                <p class="text-sm font-semibold text-foreground">{{ template.name }}</p>
                <p class="text-xs text-muted-foreground leading-tight mt-1">{{ template.description }}</p>
              </div>
            </div>
          </div>

          <!-- User Templates -->
          <div v-show="activeTemplateTab === 'user'" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div
              v-for="template in userPromptTemplates"
              :key="template.id"
              class="flex flex-col gap-3 p-4 rounded-xl border border-border bg-background hover:bg-primary/5 hover:border-primary/30 cursor-pointer transition-all duration-200 ease-in-out group hover:shadow-md hover:scale-[1.02]"
              @click="insertTemplate(template, 'user')"
              :title="template.description"
            >
              <div
                class="flex items-center justify-center h-10 w-10 rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors shadow-sm"
              >
                <component :is="template.icon" class="h-5 w-5 text-primary" />
              </div>
              <div>
                <p class="text-sm font-semibold text-foreground">{{ template.name }}</p>
                <p class="text-xs text-muted-foreground leading-tight mt-1">{{ template.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </UiScrollArea>
    </div>
  </div>
</template>
