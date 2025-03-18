<script setup lang="ts">
import { computed } from "vue";

import { useFileSystem } from "../composables/useFileSystem";

const props = defineProps<{
  filesCount: number;
  watchDir: string;
}>();

const { getFileName } = useFileSystem();

const folderName = computed(() => {
  return getFileName(props.watchDir);
});
</script>

<template>
  <footer class="status-bar">
    <div class="status-item">
      <span class="status-label">Folder:</span>
      <span class="status-value">{{ folderName }}</span>
    </div>

    <div class="status-item">
      <span class="status-label">Files:</span>
      <span class="status-value">{{ filesCount }}</span>
    </div>
  </footer>
</template>

<style scoped>
.status-bar {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  background-color: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  font-size: 0.875rem;
  color: #5f6368;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-label {
  font-weight: 500;
}
</style>
