<script setup lang="ts">
type MeResponse = { user: { email: string; sub: string } };

const me = ref<MeResponse | null>(null);
const loading = ref(false);
const err = ref<string | null>(null);

const signIn = () => window.location.assign("/auth/login");

// Call /auth/me (cookies are sent automatically on same-origin).
const fetchMe = async () => {
  loading.value = true;
  err.value = null;
  try {
    me.value = await $fetch<MeResponse>("/auth/me", {
      credentials: "same-origin", // keep relative path so Nuxt devProxy handles it
    });
  } catch (e: any) {
    if (e?.status === 401) {
      err.value = e || "Not signed in.";
    } else {
      err.value = e?.data?.message ?? e?.message ?? "Request failed";
    }
    me.value = null;
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="p-4 space-y-4">
    <div class="space-x-2">
      <button @click="signIn" class="px-3 py-1 rounded bg-black text-white">
        Sign In
      </button>
      <button @click="fetchMe" class="px-3 py-1 rounded border">
        Who am I?
      </button>
    </div>

    <div v-if="loading">Loading…</div>
    <p v-else-if="err" class="text-red-600">{{ err }}</p>
    <pre v-else-if="me" class="bg-gray-100 p-3 rounded"
      >{{ JSON.stringify(me, null, 2) }}
    </pre>
  </div>
</template>
