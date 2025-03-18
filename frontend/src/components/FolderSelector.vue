<script setup lang="ts">
import { ref } from "vue";
import { useFileSystem } from "../composables/useFileSystem";

const props = defineProps<{
  currentDir: string;
}>();

const { setWatchDir } = useFileSystem();

const isChanging = ref(false);
const newDir = ref(props.currentDir);
const error = ref<string | null>(null);

const startChanging = () => {
  isChanging.value = true;
  newDir.value = props.currentDir;
};

const cancelChanging = () => {
  isChanging.value = false;
  error.value = null;
};

const applyChanges = async () => {
  if (!newDir.value.trim()) {
    error.value = "Directory cannot be empty";
    return;
  }

  const success = await setWatchDir(newDir.value);
  if (success) {
    isChanging.value = false;
    error.value = null;
  }
};
</script>
<template>
  <div class="folder-selection">
    <div v-if="!isChanging" class="folder-display">
      <div class="folder-info">
        <h3>Watched Directory</h3>
        <div class="folder-path">{{ currentDir }}</div>
      </div>
      <button @click="startChanging" class="btn btn-change">Change</button>
    </div>

    <div v-else class="folder-edit">
      <div class="input-group">
        <label for="folder-path">Folder Path</label>
        <input
          id="folder-path"
          v-model="newDir"
          type="text"
          class="folder-input"
          placeholder="Enter folder path"
        />
        <div v-if="error" class="input-message">{{ error }}</div>
      </div>

      <div class="folder-actions">
        <button @click="cancelChanging" class="btn btn-cancel">Cancel</button>
        <button @click="applyChanges" class="btn btn-apply">Apply</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.folder-selector {
  background-color: var(--bg-secondary);
  border-radius: 4px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.folder-display {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.folder-info h3 {
  margin-bottom: 0.25rem;
  font-size: 1rem;
}

.folder-path {
  color: #5f6368;
  word-break: break-all;
}

.btn-change {
  background-color: transparent;
  color: var(--primary-color);
  border: 1px solid var(--primary-color);
}

.btn-change:hover {
  background-color: rgba(66, 133, 244, 0.05);
}

.folder-edit {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.input-group label {
  font-weight: 500;
}

.folder-input {
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 1rem;
}

.folder-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
}

.input-error {
  color: var(--error-color);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.folder-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.btn-cancel {
  background-color: transparent;
  color: #5f6368;
}

.btn-cancel:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.btn-apply {
  background-color: var(--primary-color);
  color: white;
}

.btn-apply:hover {
  background-color: #3367d6;
}
</style>
