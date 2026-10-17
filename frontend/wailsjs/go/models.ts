export namespace main {
	
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
	export class Cond {
	    kind: string;
	    not?: boolean;
	    template?: string;
	    templateScale?: number;
	    threshold?: number;
	    searchArea?: string;
	    findText?: string;
	    textExact?: boolean;
	    textNoAccent?: boolean;
	    occurrence?: number;
	    points?: Point[];
	    color?: string;
	    tolerance?: number;
	    waitMs?: number;
	    intervalMs?: number;
	    left?: string;
	    op?: string;
	    right?: string;
	    conds?: Cond[];
	
	    static createFrom(source: any = {}) {
	        return new Cond(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.not = source["not"];
	        this.template = source["template"];
	        this.templateScale = source["templateScale"];
	        this.threshold = source["threshold"];
	        this.searchArea = source["searchArea"];
	        this.findText = source["findText"];
	        this.textExact = source["textExact"];
	        this.textNoAccent = source["textNoAccent"];
	        this.occurrence = source["occurrence"];
	        this.points = this.convertValues(source["points"], Point);
	        this.color = source["color"];
	        this.tolerance = source["tolerance"];
	        this.waitMs = source["waitMs"];
	        this.intervalMs = source["intervalMs"];
	        this.left = source["left"];
	        this.op = source["op"];
	        this.right = source["right"];
	        this.conds = this.convertValues(source["conds"], Cond);
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
	export class CondPart {
	    path: number[];
	    kind: string;
	    ok: boolean;
	    known: boolean;
	    detail: string;
	    hit?: Rect;
	
	    static createFrom(source: any = {}) {
	        return new CondPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.kind = source["kind"];
	        this.ok = source["ok"];
	        this.known = source["known"];
	        this.detail = source["detail"];
	        this.hit = this.convertValues(source["hit"], Rect);
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
	export class CondTest {
	    ok: boolean;
	    known: boolean;
	    parts: CondPart[];
	    ms: number;
	
	    static createFrom(source: any = {}) {
	        return new CondTest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.known = source["known"];
	        this.parts = this.convertValues(source["parts"], CondPart);
	        this.ms = source["ms"];
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
	export class DisplayInfo {
	    index: number;
	    rect: Rect;
	    main: boolean;
	    current: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DisplayInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.rect = this.convertValues(source["rect"], Rect);
	        this.main = source["main"];
	        this.current = source["current"];
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
	export class OnFail {
	    retries?: number;
	    action: string;
	    goto?: string;
	
	    static createFrom(source: any = {}) {
	        return new OnFail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.retries = source["retries"];
	        this.action = source["action"];
	        this.goto = source["goto"];
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
	    template?: string;
	    templateScale?: number;
	    threshold?: number;
	    timeoutMs?: number;
	    intervalMs?: number;
	    searchArea?: string;
	    imageSource?: string;
	    dx?: number;
	    dy?: number;
	    findBy?: string;
	    findText?: string;
	    textExact?: boolean;
	    textNoAccent?: boolean;
	    occurrence?: number;
	    children?: Step[];
	    else?: Step[];
	    cond?: Cond;
	    name?: string;
	    repeat?: number;
	    loopMode?: string;
	    maxLoops?: number;
	    message?: string;
	    stopError?: boolean;
	    onFail?: OnFail;
	    color?: string;
	    tolerance?: number;
	
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
	        this.template = source["template"];
	        this.templateScale = source["templateScale"];
	        this.threshold = source["threshold"];
	        this.timeoutMs = source["timeoutMs"];
	        this.intervalMs = source["intervalMs"];
	        this.searchArea = source["searchArea"];
	        this.imageSource = source["imageSource"];
	        this.dx = source["dx"];
	        this.dy = source["dy"];
	        this.findBy = source["findBy"];
	        this.findText = source["findText"];
	        this.textExact = source["textExact"];
	        this.textNoAccent = source["textNoAccent"];
	        this.occurrence = source["occurrence"];
	        this.children = this.convertValues(source["children"], Step);
	        this.else = this.convertValues(source["else"], Step);
	        this.cond = this.convertValues(source["cond"], Cond);
	        this.name = source["name"];
	        this.repeat = source["repeat"];
	        this.loopMode = source["loopMode"];
	        this.maxLoops = source["maxLoops"];
	        this.message = source["message"];
	        this.stopError = source["stopError"];
	        this.onFail = this.convertValues(source["onFail"], OnFail);
	        this.color = source["color"];
	        this.tolerance = source["tolerance"];
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
	export class FindResult {
	    found: boolean;
	    score: number;
	    x: number;
	    y: number;
	    w: number;
	    h: number;
	    area: Rect;
	    ms: number;
	    read?: string[];
	
	    static createFrom(source: any = {}) {
	        return new FindResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.score = source["score"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.w = source["w"];
	        this.h = source["h"];
	        this.area = this.convertValues(source["area"], Rect);
	        this.ms = source["ms"];
	        this.read = source["read"];
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
	export class TemplateShot {
	    png: string;
	    scale: number;
	    w: number;
	    h: number;
	
	    static createFrom(source: any = {}) {
	        return new TemplateShot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.png = source["png"];
	        this.scale = source["scale"];
	        this.w = source["w"];
	        this.h = source["h"];
	    }
	}

}

