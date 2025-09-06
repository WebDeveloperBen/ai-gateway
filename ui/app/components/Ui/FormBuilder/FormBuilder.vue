<template>
  <div class="space-y-6">
    <template v-for="(field, index) in fields" :key="index">
      <template v-if="field.variant === 'Checkbox'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeCheckbox v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'Input'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeInput v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'Divider'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiDivider v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'CurrencyInput'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeCurrencyInput v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'DateField'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeDateField v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'Textarea'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeTextarea v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'FileInput'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeFileInput :name="field.name" v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'MultiSelect'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeMultiSelect v-bind="removeFields(field)">
              <template #option="{ option, isSelected }">
                <div class="flex items-center gap-2 w-full">
                  <div class="flex items-center justify-center w-4 h-4 border border-border rounded transition-colors" :class="isSelected(option) ? 'bg-primary border-primary' : 'bg-background'">
                    <svg v-if="isSelected(option)" class="w-3 h-3 text-primary-foreground" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                  <span>{{ option.label || option }}</span>
                </div>
              </template>
              <template #clear="{ clear }">
                <button 
                  class="mr-2 flex items-center justify-center hover:bg-muted/50 rounded p-1.5 transition-colors cursor-pointer z-10 relative" 
                  @click.stop="clear" 
                  type="button"
                  title="Clear all selections"
                >
                  <svg class="text-muted-foreground hover:text-foreground size-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </template>
            </UiVeeMultiSelect>
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'Select'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeSelect v-bind="removeFields(field)">
              <template v-for="(option, optIndex) in field.options" :key="optIndex">
                <option v-bind="option">{{ option.label }}</option>
              </template>
            </UiVeeSelect>
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'RadioGroup'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeRadioGroup :name="field.name" v-bind="removeFields(field)">
              <template v-for="(option, optIndex) in field.options" :key="optIndex">
                <div class="mb-2 flex items-center gap-3">
                  <UiRadioGroupItem :id="option.value" :value="option.value" />
                  <UiLabel :for="option.value">{{ option.label }}</UiLabel>
                </div>
              </template>
            </UiVeeRadioGroup>
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'PinInput'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeePinInput v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'TagsInput'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeTagsInput v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'VueformSlider'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeVueFormSlider v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
      <template v-if="field.variant === 'NativeCheckbox'">
        <slot
          v-if="field.renderIf ? field.renderIf() : true"
          :name="field.slot ? field.slot : field.name"
          v-bind="field"
        >
          <div :class="field.wrapperClass">
            <UiVeeNativeCheckbox v-bind="removeFields(field)" />
          </div>
        </slot>
      </template>
    </template>
  </div>
</template>

<script lang="ts">
import type { HTMLAttributes } from "vue"

export type FormBuilder = {
  description?: string
  hint?: string
  disabled?: boolean
  label?: string
  name: string
  placeholder?: string
  required?: boolean
  type?: string
  value?: any
  rules?: any
  class?: HTMLAttributes["class"]
  slot?: string
  wrapperClass?: HTMLAttributes["class"]
  renderIf?: () => boolean
  options?: any[]
  variant:
    | "Checkbox"
    | "NativeCheckbox"
    | "Input"
    | "Divider"
    | "CurrencyInput"
    | "DateField"
    | "FileInput"
    | "Select"
    | "Textarea"
    | "MultiSelect"
    | "PinInput"
    | "TagsInput"
    | "RadioGroup"
    | "VueformSlider"
  [key: string]: any
}
export type FormBuilderProps = {
  fields: FormBuilder[]
}
</script>

<script lang="ts" setup>
defineProps<FormBuilderProps>()

const omit = (obj: FormBuilder, keys: Array<keyof FormBuilder>) =>
  Object.fromEntries(Object.entries(obj).filter(([key]) => !keys.includes(key as keyof FormBuilder)))

const removeFields = (field: FormBuilder) => {
  const cleanedField = omit(field, ["wrapperClass", "renderIf", "variant", "slot"])

  // Transform label to formLabel only for MultiSelect (other components use label directly)
  if (field.label && field.variant === "MultiSelect") {
    cleanedField.formLabel = field.label
    delete cleanedField.label
  }

  // Add validation rules for required fields if not already provided
  if (field.required && (!("rules" in field) || !cleanedField.rules)) {
    cleanedField.rules = (value: any) => {
      if (!value || (typeof value === "string" && value.trim() === "")) {
        return `${field.label || field.name} is required`
      }
      return true
    }
  }

  // Enable validation on blur for better UX
  if (field.required && !cleanedField.validateOnMount) {
    cleanedField.validateOnMount = false // Don't validate immediately on mount
  }

  return cleanedField
}
</script>
