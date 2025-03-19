<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useFileSystem } from "./composables/useFileSystem";
import FilesList from "./components/FilesList.vue";
import FolderSelector from "./components/FolderSelector.vue";
import StatusBar from "./components/StatusBar.vue";

const {
  files,
  watchDir,
  isLoading,
  error,
  loadFiles,
  getWatchDir,
  minimizeToTray,
} = useFileSystem();

let intervalId: number;

onMounted(async () => {
  await getWatchDir();
  await loadFiles();

  // set up polling
  intervalId = setInterval(loadFiles, 5000);
});

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
});

const handleMinimize = () => {
  minimizeToTray();
};
</script>

<template>
  <div class="container-fluid min-vh-100 vh-100 p-0">
    <header class="header">
      <h1>Home Cloud</h1>

      <div class="actions">
        <button @click="handleMinimize" class="btn btn-minimmize">
          <span class="icon">_</span>
        </button>
      </div>
    </header>

    <main class="main">
      <FolderSelector :current-dir="watchDir" />

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <FilesList :files="files" :is-loading="isLoading" @refresh="loadFiles" />
    </main>

    <StatusBar
      :files-count="files.length"
      :watch-dir="watchDir"
      class="fixed-bottom"
    />
  </div>
</template>

<style>
:root {
  --primary-color: #4285f4;
  --secondary-color: #34a853;
  --error-color: #ea4335;
  --warning-color: #fbbc05;
  --text-color: #202124;
  --bg-color: #ffffff;
  --bg-secondary: #f8f9fa;
  --border-color: #dadce0;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: "Roboto", "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  color: var(--text-color);
  background-color: var(--bg-color);
}

.height-100 {
  height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: var(--primary-color);
  color: white;
}

.main {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
}

.error-message {
  background-color: #feeceb;
  color: var(--error-color);
  padding: 0.75rem;
  margin: 1rem 0;
  border-radius: 4px;
  border-left: 4px solid var(--error-color);
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn-minimize {
  background-color: transparent;
  color: white;
}

.btn-minimize:hover {
  background-color: rgba(255, 255, 255, 0.1);
}
</style>
