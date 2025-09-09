const showCreateEnvironmentModal = ref(false)

export const useGlobalModals = () => {
  const openCreateEnvironmentModal = () => {
    showCreateEnvironmentModal.value = true
  }

  const closeCreateEnvironmentModal = () => {
    showCreateEnvironmentModal.value = false
  }

  return {
    // Create Environment Modal
    showCreateEnvironmentModal, // Make it writable for v-model
    openCreateEnvironmentModal,
    closeCreateEnvironmentModal
  }
}

