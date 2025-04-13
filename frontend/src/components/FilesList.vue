<script setup lang="ts">
import { computed } from "vue";
import type { FileInfo } from "../types";
import { useFileSystem } from "../composables/useFileSystem";

const props = defineProps<{
  files: FileInfo[];
  isLoading: boolean;
}>();

const emits = defineEmits<{
  (e: "refresh"): void;
}>();

const { getFileName, formatFileSize } = useFileSystem();

const sortedFiles = computed(() => {
  return [...props.files].sort((a, b) => {
    // Sort by status first (NOT_SYNCED first, then SYNCING, etc.)
    const statusOrder = {
      NOT_SYNCED: 0,
      SYNCING: 1,
      SYNCED: 2,
      ERROR: 3,
    };

    const statusA = statusOrder[a.Status as keyof typeof statusOrder];
    const statusB = statusOrder[b.Status as keyof typeof statusOrder];

    if (statusA !== statusB) {
      return statusA - statusB;
    }

    // Sort by name
    return getFileName(a.Path).localeCompare(getFileName(b.Path));
  });
});

const getStatusClass = (status: string) => {
  switch (status) {
    case "SYNCED":
      return "status-synced";
    case "SYNCING":
      return "status-syncing";
    case "ERROR":
      return "status-error";
    default:
      return "status-not-synced";
  }
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case "SYNCED":
      return "✅";
    case "SYNCING":
      return "🔄";
    case "ERROR":
      return "❌";
    default:
      return "⏳";
  }
};
</script>

<template>
  <div class="files-container">
    <div class="files-header">
      <h2>Files</h2>
      <button
        @click="emits('refresh')"
        class="btn btn-refresh"
        :disabled="isLoading"
      >
        Refresh
      </button>
    </div>

    <div v-if="isLoading && files.length === 0" class="loading-indicator">
      Loading files...
    </div>

    <div v-else-if="files.length === 0" class="empty-state">
      No files found in the watched directory
    </div>

    <div v-else class="files-list">
      <div class="files-list-header">
        <span class="file-name-header">Name</span>
        <span class="file-status-header">Status</span>
        <span class="file-size-header">Size</span>
      </div>

      <div id="filesAccordion" class="accordion">
        <div v-for="file in sortedFiles" :key="file.Path" class="file-item">
          <div v-if="file.IsDirectory" class="file-name accordion-item">
            <span class="accordion-header">{{ getFileName(file.Path) }}</span>
            <div
              class="accordion-collapse collapse"
              data-bs-parent="filesAccordion"
            >
              <!-- <FilesList
                :files="file.FilesContent"
                :is-loading="isLoading"
                @refresh="emits('refresh')"
              /> -->
            </div>
          </div>
          <div v-else>
            <div class="file-name">{{ getFileName(file.Path) }}</div>
            <div class="file-status" :class="getStatusClass(file.Status)">
              <span class="status-icon">{{ getStatusIcon(file.Status) }}</span>
              {{ file.Status.replace("_", " ") }}
            </div>
            <div class="file-size">{{ formatFileSize(file.Size) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.files-container {
  margin-top: 1rem;
}

.files-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.btn-refresh {
  background-color: var(--primary-color);
  color: white;
}

.btn-refresh:hover {
  background-color: #3367d6;
}

.btn-refresh:disabled {
  background-color: #a4c2f4;
  cursor: not-allowed;
}

.loading-indicator,
.empty-state {
  padding: 2rem;
  text-align: center;
  color: #5f6368;
  background-color: var(--bg-secondary);
  border-radius: 4px;
}

.files-list {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.files-list-header {
  display: grid;
  grid-template-columns: 1fr 180px 100px;
  background-color: var(--bg-secondary);
  padding: 0.75rem 1rem;
  font-weight: 500;
  border-bottom: 1px solid var(--border-color);
}

.file-item {
  display: grid;
  grid-template-columns: 1fr 180px 100px;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
  align-items: center;
}

.file-item:last-child {
  border-bottom: none;
}

.file-item:hover {
  background-color: #f1f3f4;
}

.file-name {
  color: black;
}

.file-status {
  display: flex;
  align-items: center;
  font-weight: 500;
}

.status-icon {
  margin-right: 0.5rem;
}

.status-synced {
  color: var(--secondary-color);
}

.status-syncing {
  color: var(--primary-color);
}

.status-error {
  color: var(--error-color);
}

.status-not-synced {
  color: var(--warning-color);
}

.file-size {
  text-align: right;
  color: #5f6368;
}
</style>
