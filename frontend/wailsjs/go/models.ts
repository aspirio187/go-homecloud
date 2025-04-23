export namespace models {
	
	export class FileInfo {
	    path: string;
	    status: string;
	    // Go type: time
	    lastModified: any;
	    size: number;
	    isDownloaded: boolean;
	    isDirectory: boolean;
	    version: number;
	    checksum?: string;
	    // Go type: time
	    lastSynced: any;
	    filesContent?: Record<string, FileInfo>;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.status = source["status"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.size = source["size"];
	        this.isDownloaded = source["isDownloaded"];
	        this.isDirectory = source["isDirectory"];
	        this.version = source["version"];
	        this.checksum = source["checksum"];
	        this.lastSynced = this.convertValues(source["lastSynced"], null);
	        this.filesContent = this.convertValues(source["filesContent"], FileInfo, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

