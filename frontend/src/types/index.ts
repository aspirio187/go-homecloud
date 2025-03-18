export interface FileInfo {
  Path: string;
  Status: "NOT_SYNCED" | "SYNCING" | "SYNCED" | "ERROR";
  LastModified: string;
  Size: number;
  IsDownloaded: boolean;
}

export interface SyncState {
  watchDir: string;
  files: FileInfo[];
  isLoading: boolean;
  error: string | null;
}
