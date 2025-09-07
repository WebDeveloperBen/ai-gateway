<script setup lang="ts">
import { Settings } from "lucide-vue-next"

const { modelState, modalState, modalActions } = usePlaygroundState()!

const isOpen = computed({
  get: () => modalState.showParametersModal,
  set: (value) => {
    if (value) {
      modalActions.openParametersModal()
    } else {
      modalActions.closeParametersModal()
    }
  }
})

function closeModal() {
  modalActions.closeParametersModal()
}

function resetToDefaults() {
  modelState.parameters = {
    temperature: 0.7,
    maxTokens: 1000,
    topP: 1.0,
    frequencyPenalty: 0.0,
    presencePenalty: 0.0
  }
}

function applySettings() {
  closeModal()
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="min-w-3xl">
      <template #header>
        <div>
          <h3 class="font-semibold text-lg flex items-center gap-2">
            <Settings class="h-5 w-5 text-primary" />
            Model Parameters
          </h3>
          <p class="text-sm text-muted-foreground mt-1">Fine-tune AI behavior and output</p>
        </div>
      </template>

      <template #content>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Temperature -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <label class="text-sm font-medium text-foreground">Temperature</label>
              <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">
                {{ modelState.parameters.temperature }}
              </span>
            </div>
            <input
              type="range"
              v-model.number="modelState.parameters.temperature"
              min="0"
              max="2"
              step="0.1"
              class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
            />
            <p class="text-xs text-muted-foreground">
              Controls randomness. Lower = more focused, Higher = more creative
            </p>
          </div>

          <!-- Max Tokens -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <label class="text-sm font-medium text-foreground">Max Tokens</label>
              <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">
                {{ modelState.parameters.maxTokens }}
              </span>
            </div>
            <input
              type="range"
              v-model.number="modelState.parameters.maxTokens"
              min="50"
              :max="modelState.selectedModelData?.maxTokens || 4000"
              step="50"
              class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
            />
            <p class="text-xs text-muted-foreground">Maximum response length</p>
          </div>

          <!-- Top P -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <label class="text-sm font-medium text-foreground">Top P</label>
              <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">
                {{ modelState.parameters.topP }}
              </span>
            </div>
            <input
              type="range"
              v-model.number="modelState.parameters.topP"
              min="0.1"
              max="1"
              step="0.05"
              class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
            />
            <p class="text-xs text-muted-foreground">Controls diversity via nucleus sampling</p>
          </div>

          <!-- Frequency Penalty -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <label class="text-sm font-medium text-foreground">Frequency Penalty</label>
              <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">
                {{ modelState.parameters.frequencyPenalty }}
              </span>
            </div>
            <input
              type="range"
              v-model.number="modelState.parameters.frequencyPenalty"
              min="-2"
              max="2"
              step="0.1"
              class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
            />
            <p class="text-xs text-muted-foreground">Reduces repetition of frequent tokens</p>
          </div>

          <!-- Presence Penalty -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <label class="text-sm font-medium text-foreground">Presence Penalty</label>
              <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">
                {{ modelState.parameters.presencePenalty }}
              </span>
            </div>
            <input
              type="range"
              v-model.number="modelState.parameters.presencePenalty"
              min="-2"
              max="2"
              step="0.1"
              class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
            />
            <p class="text-xs text-muted-foreground">Encourages talking about new topics</p>
          </div>

          <!-- Current Model Info -->
          <div class="md:col-span-2 p-4 bg-muted/20 rounded-lg border">
            <div v-if="modelState.selectedModelData" class="text-sm">
              <h4 class="font-medium mb-2">Current Model: {{ modelState.selectedModelData.name }}</h4>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs text-muted-foreground">
                <div>
                  <span class="font-medium">Max Tokens:</span>
                  {{ modelState.selectedModelData.maxTokens.toLocaleString() }}
                </div>
                <div>
                  <span class="font-medium">Input Cost:</span>
                  {{ formatCurrency(modelState.selectedModelData.costPer1kTokens.input) }}/1k
                </div>
                <div>
                  <span class="font-medium">Output Cost:</span>
                  {{ formatCurrency(modelState.selectedModelData.costPer1kTokens.output) }}/1k
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #footer>
        <UiDialogFooter class="flex gap-3 justify-between pt-6 border-t mt-6">
          <UiButton variant="outline" @click="resetToDefaults"> Reset to Defaults </UiButton>

          <div class="flex gap-3">
            <UiButton variant="outline" @click="closeModal">Cancel</UiButton>
            <UiButton @click="applySettings">Apply Settings</UiButton>
          </div>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
