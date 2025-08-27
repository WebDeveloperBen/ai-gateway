<script lang="ts" setup>
import { promiseTimeout } from "@vueuse/core";
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue";
import { toast } from "vue-sonner";
import * as z from "zod";

const schema = z.object({
  firstName: z.string().min(2),
  lastName: z.string().min(2),
  email: z.email(),
  password: z.string().min(6),
});

const { handleSubmit, isSubmitting, resetForm } = useForm<
  z.infer<typeof schema>
>({
  name: "form-builder-full",
  validationSchema: toTypedSchema(schema),
  initialValues: {
    firstName: "Jane",
  },
});

const submit = handleSubmit((values: any) => {
  try {
    console.log(values);
    toast.success("Form submitted successfully", {
      description: "We will reset the form in 5 seconds",
    });
    promiseTimeout(5000).then(() => {
      resetForm();
    });
  } catch (error) {
    console.log(error);
    toast.error("Form submission failed", {
      description: "Please try again",
    });
  }
});

const form: FormBuilder[] = [
  {
    variant: "Input",
    label: "First Name",
    name: "firstName",
    placeholder: "Enter your first name",
    required: true,
    hint: "Enter your first name",
    wrapperClass: tw`col-span-full md:col-span-6`,
  },
  {
    variant: "Input",
    label: "Last Name",
    name: "lastName",
    placeholder: "Enter your last name",
    required: true,
    hint: "Enter your last name",
    wrapperClass: tw`col-span-full md:col-span-6`,
  },
  {
    variant: "Input",
    label: "Email",
    name: "email",
    placeholder: "igugowuj@jakomka.edu",
    type: "email",
    required: true,
    icon: "lucide:mail",
    wrapperClass: tw`col-span-full md:col-span-6`,
  },
  {
    variant: "Input",
    label: "Password",
    name: "password",
    placeholder: "••••••••",
    icon: "lucide:lock",
    required: true,
    type: "password",
    wrapperClass: tw`col-span-full md:col-span-6`,
  },
];
</script>
<template>
  <form @submit="submit">
    <fieldset :disabled="isSubmitting">
      <UiFormBuilder class="grid grid-cols-12 gap-5" :fields="form" />
    </fieldset>
  </form>
</template>
