export namespace sync {
	
	export class FileInfo {
	    Path: string;
	    Status: string;
	    // Go type: time
	    LastModified: any;
	    Size: number;
	    IsDownloaded: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Status = source["Status"];
	        this.LastModified = this.convertValues(source["LastModified"], null);
	        this.Size = source["Size"];
	        this.IsDownloaded = source["IsDownloaded"];
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

