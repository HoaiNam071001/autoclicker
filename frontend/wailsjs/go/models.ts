export namespace main {
	
	export class InputTarget {
	    stepId: string;
	    field: string;
	
	    static createFrom(source: any = {}) {
	        return new InputTarget(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stepId = source["stepId"];
	        this.field = source["field"];
	    }
	}
	export class ScriptInput {
	    key: string;
	    label: string;
	    type: string;
	    default?: any;
	    required?: boolean;
	    min?: number;
	    max?: number;
	    options?: string[];
	    hint?: string;
	    targets?: InputTarget[];
	
	    static createFrom(source: any = {}) {
	        return new ScriptInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.default = source["default"];
	        this.required = source["required"];
	        this.min = source["min"];
	        this.max = source["max"];
	        this.options = source["options"];
	        this.hint = source["hint"];
	        this.targets = this.convertValues(source["targets"], InputTarget);
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
	export class Point {
	    x: number;
	    y: number;
	    ax?: string;
	    ay?: string;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.ax = source["ax"];
	        this.ay = source["ay"];
	    }
	}
	export class Step {
	    id: string;
	    type: string;
	    points: Point[];
	    button: string;
	    clicks: number;
	    clickInterval: number;
	    holdMs: number;
	    dragMs: number;
	    dragFromPrev: boolean;
	    dragRelative: boolean;
	    text: string;
	    keys: string[];
	    keyMode: string;
	    modifiers: string[];
	    delayAfter: number;
	    jitterMs: number;
	    jitterPx: number;
	    enabled: boolean;
	    direction?: string;
	    amount?: number;
	    target?: string;
	    app?: string;
	
	    static createFrom(source: any = {}) {
	        return new Step(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.points = this.convertValues(source["points"], Point);
	        this.button = source["button"];
	        this.clicks = source["clicks"];
	        this.clickInterval = source["clickInterval"];
	        this.holdMs = source["holdMs"];
	        this.dragMs = source["dragMs"];
	        this.dragFromPrev = source["dragFromPrev"];
	        this.dragRelative = source["dragRelative"];
	        this.text = source["text"];
	        this.keys = source["keys"];
	        this.keyMode = source["keyMode"];
	        this.modifiers = source["modifiers"];
	        this.delayAfter = source["delayAfter"];
	        this.jitterMs = source["jitterMs"];
	        this.jitterPx = source["jitterPx"];
	        this.enabled = source["enabled"];
	        this.direction = source["direction"];
	        this.amount = source["amount"];
	        this.target = source["target"];
	        this.app = source["app"];
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
	export class Rect {
	    x: number;
	    y: number;
	    w: number;
	    h: number;
	
	    static createFrom(source: any = {}) {
	        return new Rect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.w = source["w"];
	        this.h = source["h"];
	    }
	}
	export class Script {
	    version: number;
	    name: string;
	    region?: Rect;
	    loops: number;
	    maxMinutes: number;
	    steps: Step[];
	    anchorMode?: string;
	    inputs?: ScriptInput[];
	
	    static createFrom(source: any = {}) {
	        return new Script(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.name = source["name"];
	        this.region = this.convertValues(source["region"], Rect);
	        this.loops = source["loops"];
	        this.maxMinutes = source["maxMinutes"];
	        this.steps = this.convertValues(source["steps"], Step);
	        this.anchorMode = source["anchorMode"];
	        this.inputs = this.convertValues(source["inputs"], ScriptInput);
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
	export class Draft {
	    script: Script;
	    path: string;
	    dirty: boolean;
	    savedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new Draft(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.script = this.convertValues(source["script"], Script);
	        this.path = source["path"];
	        this.dirty = source["dirty"];
	        this.savedAt = source["savedAt"];
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
	export class InputPreset {
	    name: string;
	    values: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new InputPreset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.values = source["values"];
	    }
	}
	export class InputMemory {
	    last?: Record<string, any>;
	    lastPreset?: string;
	    presets?: InputPreset[];
	    autoRun?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InputMemory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.last = source["last"];
	        this.lastPreset = source["lastPreset"];
	        this.presets = this.convertValues(source["presets"], InputPreset);
	        this.autoRun = source["autoRun"];
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
	
	
	
	
	
	export class ScriptFile {
	    script: Script;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ScriptFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.script = this.convertValues(source["script"], Script);
	        this.path = source["path"];
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
	
	export class Settings {
	    theme: string;
	    recent?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.recent = source["recent"];
	    }
	}
	
	export class StepTypeDoc {
	    type: string;
	    label: string;
	    desc: string;
	    fields: string[];
	    pinnable: string[];
	    example: string;
	
	    static createFrom(source: any = {}) {
	        return new StepTypeDoc(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.label = source["label"];
	        this.desc = source["desc"];
	        this.fields = source["fields"];
	        this.pinnable = source["pinnable"];
	        this.example = source["example"];
	    }
	}

}

