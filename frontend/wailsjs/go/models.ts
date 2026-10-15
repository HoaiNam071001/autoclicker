export namespace main {
	
	export class SpikeResult {
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new SpikeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}

}

