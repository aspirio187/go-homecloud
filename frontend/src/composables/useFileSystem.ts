import { ref, computed } from "vue";
import type { FileInfo } from "../types";

export function useFileSystem() {
  const files = ref<FileInfo[]>([]);
  const watchDir = ref<string>("");
  const isLoading = ref<boolean>(false);
  const error = ref<string | null>(null);

  const syncedFile = computed(() =>
    files.value.filter((file) => file.Status === "SYNCED")
  );

  const pendingFiles = computed(() =>
    files.value.filter(
      (file) => file.Status === "NOT_SYNCED" || file.Status === "SYNCING"
    )
  );

  const getFileName = (path: string): string => {
    return path.split(/[\/\\]/).pop() || path;
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + " KB";
    if (bytes < 1024 * 1024 * 1024)
      return (bytes / (1024 * 1024)).toFixed(2) + " MB";
    return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
  };

  const loadFiles = async () => {
    isLoading.value = true;
    error.value = null;

    try {
      files.value = await window.go.main.App.GetFiles();
    } catch (e) {
      console.error("Failed to load files:", e);
      error.value = "Failed to load files";
    } finally {
      isLoading.value = false;
    }
  };

  const getWatchDir = async () => {
    try {
      watchDir.value = await window.go.main.App.GetWatchDir();
    } catch (e) {
      console.error("Failed to get watch dir:", e);
      error.value = "Failed to get watch dir";
    }
  };

  const setWatchDir = async (dir: string) => {
    isLoading.value = true;
    error.value = null;

    try {
      await window.go.main.App.SetWatchDir(dir);
      watchDir.value = dir;
      await loadFiles();
      return true;
    } catch (e) {
      console.error("Failed to set watch directory:", e);
      error.value = "Failed to set watch directory";
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  const minimizeToTray = async () => {
    try {
      await window.go.main.App.MinimizeToTray;
    } catch (e) {
      console.error("Failed to minimize to tray:", e);
    }
  };

  return {
    files,
    syncedFile,
    pendingFiles,
    watchDir,
    isLoading,
    error,
    getFileName,
    formatFileSize,
    loadFiles,
    getWatchDir,
    setWatchDir,
    minimizeToTray,
  };
}
