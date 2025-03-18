export {};

declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetFiles(): Promise<any[]>;
          GetWatchDir(): Promise<string>;
          SetWatchDir(dir: string): Promise<void>;
          MinimizeToTray(): Promise<void>;
        };
      };
    };
  }
}
