export namespace core {
	
	export class ChineseFunction {
	    chineseName: string;
	    goFunction: string;
	    description: string;
	    category: string;
	    params: string[];
	    returnType: string;
	
	    static createFrom(source: any = {}) {
	        return new ChineseFunction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chineseName = source["chineseName"];
	        this.goFunction = source["goFunction"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.params = source["params"];
	        this.returnType = source["returnType"];
	    }
	}
	export class CodeRowData {
	    event: string;
	    code: string;
	
	    static createFrom(source: any = {}) {
	        return new CodeRowData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.event = source["event"];
	        this.code = source["code"];
	    }
	}
	export class CodeTemplate {
	    name: string;
	    description: string;
	    category: string;
	    code: string;
	
	    static createFrom(source: any = {}) {
	        return new CodeTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.code = source["code"];
	    }
	}

}

export namespace models {
	
	export class Project {
	    id: number;
	    name: string;
	    description: string;
	    path: string;
	    type: string;
	    author: string;
	    version: string;
	    icon: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    // Go type: time
	    lastOpened: any;
	    isFavorite: boolean;
	    tags: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.path = source["path"];
	        this.type = source["type"];
	        this.author = source["author"];
	        this.version = source["version"];
	        this.icon = source["icon"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.lastOpened = this.convertValues(source["lastOpened"], null);
	        this.isFavorite = source["isFavorite"];
	        this.tags = source["tags"];
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
	export class ProjectTemplate {
	    id: number;
	    name: string;
	    description: string;
	    category: string;
	    icon: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.icon = source["icon"];
	        this.content = source["content"];
	    }
	}

}

export namespace services {
	
	export class AIConfig {
	    provider: string;
	    localModelPath: string;
	    localModelName: string;
	    ollamaUrl: string;
	    ollamaModel: string;
	    cloudApiKey: string;
	    cloudApiUrl: string;
	    cloudModel: string;
	    enableGPU: boolean;
	    gpuLayers: number;
	    maxTokens: number;
	    temperature: number;
	    topP: number;
	    contextSize: number;
	    enableCodeCompletion: boolean;
	    completionDelay: number;
	    enableChatHistory: boolean;
	    systemPrompt: string;
	
	    static createFrom(source: any = {}) {
	        return new AIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.localModelPath = source["localModelPath"];
	        this.localModelName = source["localModelName"];
	        this.ollamaUrl = source["ollamaUrl"];
	        this.ollamaModel = source["ollamaModel"];
	        this.cloudApiKey = source["cloudApiKey"];
	        this.cloudApiUrl = source["cloudApiUrl"];
	        this.cloudModel = source["cloudModel"];
	        this.enableGPU = source["enableGPU"];
	        this.gpuLayers = source["gpuLayers"];
	        this.maxTokens = source["maxTokens"];
	        this.temperature = source["temperature"];
	        this.topP = source["topP"];
	        this.contextSize = source["contextSize"];
	        this.enableCodeCompletion = source["enableCodeCompletion"];
	        this.completionDelay = source["completionDelay"];
	        this.enableChatHistory = source["enableChatHistory"];
	        this.systemPrompt = source["systemPrompt"];
	    }
	}
	export class BuildConfig {
	    appName: string;
	    version: string;
	    description: string;
	    company: string;
	    copyright: string;
	    iconPath: string;
	    outputPath: string;
	    targetOS: string;
	    targetArch: string;
	    optimize: boolean;
	    stripDebug: boolean;
	    upxCompress: boolean;
	    embedAssets: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BuildConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.company = source["company"];
	        this.copyright = source["copyright"];
	        this.iconPath = source["iconPath"];
	        this.outputPath = source["outputPath"];
	        this.targetOS = source["targetOS"];
	        this.targetArch = source["targetArch"];
	        this.optimize = source["optimize"];
	        this.stripDebug = source["stripDebug"];
	        this.upxCompress = source["upxCompress"];
	        this.embedAssets = source["embedAssets"];
	    }
	}
	export class BuildResult {
	    success: boolean;
	    message: string;
	    output: string;
	    exePath: string;
	    progress: number;
	
	    static createFrom(source: any = {}) {
	        return new BuildResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.output = source["output"];
	        this.exePath = source["exePath"];
	        this.progress = source["progress"];
	    }
	}
	export class ChatMessage {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class ChatResponse {
	    content: string;
	    done: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.done = source["done"];
	        this.error = source["error"];
	    }
	}
	export class CompletionResponse {
	    text: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CompletionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.error = source["error"];
	    }
	}

}

