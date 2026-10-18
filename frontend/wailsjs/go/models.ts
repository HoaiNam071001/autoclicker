export namespace engine {
	
	export class CondPart {
	    path: number[];
	    case?: string;
	    kind: string;
	    ok: boolean;
	    known: boolean;
	    detail: string;
	    hit?: model.Rect;
	    score?: number;
	    text?: string;
	
	    static createFrom(source: any = {}) {
	        return new CondPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.case = source["case"];
	        this.kind = source["kind"];
	        this.ok = source["ok"];
	        this.known = source["known"];
	        this.detail = source["detail"];
	        this.hit = this.convertValues(source["hit"], model.Rect);
	        this.score = source["score"];
	        this.text = source["text"];
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
	    area: model.Rect;
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
	        this.area = this.convertValues(source["area"], model.Rect);
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
	export class NodeTest {
	    port: string;
	    known: boolean;
	    parts: CondPart[];
	    value?: string;
	    ms: number;
	
	    static createFrom(source: any = {}) {
	        return new NodeTest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.known = source["known"];
	        this.parts = this.convertValues(source["parts"], CondPart);
	        this.value = source["value"];
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
	export class WebFindResult {
	    found: boolean;
	    x: number;
	    y: number;
	    w: number;
	    h: number;
	    text: string;
	    selector: string;
	    covered: string;
	    offView: boolean;
	    hidden: number;
	    tab: string;
	    noTab: boolean;
	    ms: number;
	
	    static createFrom(source: any = {}) {
	        return new WebFindResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.w = source["w"];
	        this.h = source["h"];
	        this.text = source["text"];
	        this.selector = source["selector"];
	        this.covered = source["covered"];
	        this.offView = source["offView"];
	        this.hidden = source["hidden"];
	        this.tab = source["tab"];
	        this.noTab = source["noTab"];
	        this.ms = source["ms"];
	    }
	}

}

export namespace main {
	
	export class DisplayInfo {
	    index: number;
	    rect: model.Rect;
	    main: boolean;
	    current: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DisplayInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.rect = this.convertValues(source["rect"], model.Rect);
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
	export class RunOptions {
	    countdown: number;
	    path?: string;
	    source?: string;
	    from?: string;
	    step?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RunOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.countdown = source["countdown"];
	        this.path = source["path"];
	        this.source = source["source"];
	        this.from = source["from"];
	        this.step = source["step"];
	    }
	}
	export class ScriptFile {
	    script: model.Script;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ScriptFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.script = this.convertValues(source["script"], model.Script);
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

export namespace model {
	
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
	export class Case {
	    key: string;
	    label?: string;
	    op?: string;
	    match?: string;
	    cond?: Cond;
	
	    static createFrom(source: any = {}) {
	        return new Case(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.op = source["op"];
	        this.match = source["match"];
	        this.cond = this.convertValues(source["cond"], Cond);
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
	
	export class Edge {
	    id: string;
	    from: string;
	    port: string;
	    to: string;
	
	    static createFrom(source: any = {}) {
	        return new Edge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.from = source["from"];
	        this.port = source["port"];
	        this.to = source["to"];
	    }
	}
	export class FlowGuideDoc {
	    flow: string;
	    commonFields: string;
	    example: string;
	
	    static createFrom(source: any = {}) {
	        return new FlowGuideDoc(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.flow = source["flow"];
	        this.commonFields = source["commonFields"];
	        this.example = source["example"];
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
	export class LogLine {
	    atMs: number;
	    level: string;
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new LogLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.atMs = source["atMs"];
	        this.level = source["level"];
	        this.msg = source["msg"];
	    }
	}
	
	export class Pos {
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new Pos(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
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
	export class Retry {
	    count?: number;
	    delayMs?: number;
	
	    static createFrom(source: any = {}) {
	        return new Retry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.delayMs = source["delayMs"];
	    }
	}
	export class StepError {
	    code: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new StepError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}
	export class StepResult {
	    nodeId: string;
	    ref: string;
	    status: string;
	    port?: string;
	    output?: Record<string, any>;
	    error?: StepError;
	    logs?: LogLine[];
	    attempts: number;
	    loop: number;
	    shot?: string;
	    // Go type: time
	    startedAt: any;
	    durationMs: number;
	
	    static createFrom(source: any = {}) {
	        return new StepResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nodeId = source["nodeId"];
	        this.ref = source["ref"];
	        this.status = source["status"];
	        this.port = source["port"];
	        this.output = source["output"];
	        this.error = this.convertValues(source["error"], StepError);
	        this.logs = this.convertValues(source["logs"], LogLine);
	        this.attempts = source["attempts"];
	        this.loop = source["loop"];
	        this.shot = source["shot"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.durationMs = source["durationMs"];
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
	export class RunRecord {
	    id: string;
	    script: string;
	    source: string;
	    // Go type: time
	    startedAt: any;
	    durationMs: number;
	    status: string;
	    error?: string;
	    message?: string;
	    nodes: number;
	    inputs?: Record<string, any>;
	    results: StepResult[];
	    truncated?: number;
	
	    static createFrom(source: any = {}) {
	        return new RunRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.script = source["script"];
	        this.source = source["source"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.durationMs = source["durationMs"];
	        this.status = source["status"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.nodes = source["nodes"];
	        this.inputs = source["inputs"];
	        this.results = this.convertValues(source["results"], StepResult);
	        this.truncated = source["truncated"];
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
	export class RunSummary {
	    id: string;
	    script: string;
	    source: string;
	    // Go type: time
	    startedAt: any;
	    durationMs: number;
	    status: string;
	    error?: string;
	    message?: string;
	    nodes: number;
	
	    static createFrom(source: any = {}) {
	        return new RunSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.script = source["script"];
	        this.source = source["source"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.durationMs = source["durationMs"];
	        this.status = source["status"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.nodes = source["nodes"];
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
	export class Step {
	    id: string;
	    type: string;
	    ref: string;
	    name?: string;
	    pos?: Pos;
	    breakpoint?: boolean;
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
	    webBrowser?: boolean;
	    selectors?: string[];
	    tab?: string;
	    webGone?: boolean;
	    template?: string;
	    templateScale?: number;
	    threshold?: number;
	    timeoutMs?: number;
	    intervalMs?: number;
	    searchArea?: string;
	    imageSource?: string;
	    imageFrom?: string;
	    dx?: number;
	    dy?: number;
	    findBy?: string;
	    findText?: string;
	    textExact?: boolean;
	    textNoAccent?: boolean;
	    occurrence?: number;
	    retry?: Retry;
	    color?: string;
	    tolerance?: number;
	    cond?: Cond;
	    cases?: Case[];
	    value?: string;
	    max?: number;
	    endStatus?: string;
	    message?: string;
	    varName?: string;
	    varOp?: string;
	
	    static createFrom(source: any = {}) {
	        return new Step(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.ref = source["ref"];
	        this.name = source["name"];
	        this.pos = this.convertValues(source["pos"], Pos);
	        this.breakpoint = source["breakpoint"];
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
	        this.webBrowser = source["webBrowser"];
	        this.selectors = source["selectors"];
	        this.tab = source["tab"];
	        this.webGone = source["webGone"];
	        this.template = source["template"];
	        this.templateScale = source["templateScale"];
	        this.threshold = source["threshold"];
	        this.timeoutMs = source["timeoutMs"];
	        this.intervalMs = source["intervalMs"];
	        this.searchArea = source["searchArea"];
	        this.imageSource = source["imageSource"];
	        this.imageFrom = source["imageFrom"];
	        this.dx = source["dx"];
	        this.dy = source["dy"];
	        this.findBy = source["findBy"];
	        this.findText = source["findText"];
	        this.textExact = source["textExact"];
	        this.textNoAccent = source["textNoAccent"];
	        this.occurrence = source["occurrence"];
	        this.retry = this.convertValues(source["retry"], Retry);
	        this.color = source["color"];
	        this.tolerance = source["tolerance"];
	        this.cond = this.convertValues(source["cond"], Cond);
	        this.cases = this.convertValues(source["cases"], Case);
	        this.value = source["value"];
	        this.max = source["max"];
	        this.endStatus = source["endStatus"];
	        this.message = source["message"];
	        this.varName = source["varName"];
	        this.varOp = source["varOp"];
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
	    start: string;
	    edges: Edge[];
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
	        this.start = source["start"];
	        this.edges = this.convertValues(source["edges"], Edge);
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
	
	
	
	
	export class StepTypeDoc {
	    type: string;
	    label: string;
	    desc: string;
	    fields: string[];
	    pinnable: string[];
	    example: string;
	    outputs: string[];
	    ports: string[];
	
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
	        this.outputs = source["outputs"];
	        this.ports = source["ports"];
	    }
	}

}

export namespace samples {
	
	export class Info {
	    id: string;
	    name: string;
	    description: string;
	    app: string;
	    level: number;
	    page: string;
	    real: boolean;
	    prep?: string;
	    steps: string[];
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.app = source["app"];
	        this.level = source["level"];
	        this.page = source["page"];
	        this.real = source["real"];
	        this.prep = source["prep"];
	        this.steps = source["steps"];
	    }
	}

}

export namespace store {
	
	export class Draft {
	    script: model.Script;
	    path: string;
	    dirty: boolean;
	    savedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new Draft(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.script = this.convertValues(source["script"], model.Script);
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

}

export namespace web {
	
	export class Picked {
	    selectors: string[];
	    text: string;
	    tab: string;
	
	    static createFrom(source: any = {}) {
	        return new Picked(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selectors = source["selectors"];
	        this.text = source["text"];
	        this.tab = source["tab"];
	    }
	}
	export class Tab {
	    title: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new Tab(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.url = source["url"];
	    }
	}
	export class Status {
	    open: boolean;
	    tabs: Tab[];
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.open = source["open"];
	        this.tabs = this.convertValues(source["tabs"], Tab);
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

