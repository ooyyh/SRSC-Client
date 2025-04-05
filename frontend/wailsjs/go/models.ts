export namespace main {
	
	export class ObjectInfo {
	    key: string;
	    size: number;
	    // Go type: time
	    lastModified: any;
	    etag: string;
	    contentType: string;
	    storageClass: string;
	    metadata: {[key: string]: string};
	    versionId: string;
	
	    static createFrom(source: any = {}) {
	        return new ObjectInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.etag = source["etag"];
	        this.contentType = source["contentType"];
	        this.storageClass = source["storageClass"];
	        this.metadata = source["metadata"];
	        this.versionId = source["versionId"];
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
	export class S3Manager {
	
	
	    static createFrom(source: any = {}) {
	        return new S3Manager(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace nodes {
	
	export class BucketInfo {
	    name: string;
	    // Go type: time
	    creationDate: any;
	    usedSpace: number;
	    totalObjects: number;
	    versioningEnabled: boolean;
	    publicAccessBlocked: boolean;
	    hasPolicy: boolean;
	    encryptionEnabled: boolean;
	    encryptionType: string;
	    hasLifecycleRules: boolean;
	    lifecycleRulesCount: number;
	    region: string;
	    websiteEnabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BucketInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.creationDate = this.convertValues(source["creationDate"], null);
	        this.usedSpace = source["usedSpace"];
	        this.totalObjects = source["totalObjects"];
	        this.versioningEnabled = source["versioningEnabled"];
	        this.publicAccessBlocked = source["publicAccessBlocked"];
	        this.hasPolicy = source["hasPolicy"];
	        this.encryptionEnabled = source["encryptionEnabled"];
	        this.encryptionType = source["encryptionType"];
	        this.hasLifecycleRules = source["hasLifecycleRules"];
	        this.lifecycleRulesCount = source["lifecycleRulesCount"];
	        this.region = source["region"];
	        this.websiteEnabled = source["websiteEnabled"];
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
	export class Node {
	    NodeName: string;
	    EndPoint: string;
	    AccessKey: string;
	    SecretKey: string;
	    Region: string;
	
	    static createFrom(source: any = {}) {
	        return new Node(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.NodeName = source["NodeName"];
	        this.EndPoint = source["EndPoint"];
	        this.AccessKey = source["AccessKey"];
	        this.SecretKey = source["SecretKey"];
	        this.Region = source["Region"];
	    }
	}
	export class NodeBucketInfo {
	    nodeName: string;
	    endPoint: string;
	    buckets: BucketInfo[];
	
	    static createFrom(source: any = {}) {
	        return new NodeBucketInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nodeName = source["nodeName"];
	        this.endPoint = source["endPoint"];
	        this.buckets = this.convertValues(source["buckets"], BucketInfo);
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

