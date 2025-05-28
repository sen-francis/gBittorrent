export namespace services {
	
	export class File {
	    file: string;
	    error: any;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.error = source["error"];
	    }
	}

}

