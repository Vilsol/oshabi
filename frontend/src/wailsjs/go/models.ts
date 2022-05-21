export namespace main {
	
	export class ConvertedConfig {
	    market_prices: {[key: string]: string};
	    prices: {[key: string]: string};
	    language: string;
	    league: string;
	    messages: {[key: string]: string};
	    name: string;
	    stream: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConvertedConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_prices = source["market_prices"];
	        this.prices = source["prices"];
	        this.language = source["language"];
	        this.league = source["league"];
	        this.messages = source["messages"];
	        this.name = source["name"];
	        this.stream = source["stream"];
	    }
	}

}

export namespace types {
	
	export class ParsedListing {
	    type: string;
	    count: number;
	    level: number;
	
	    static createFrom(source: any = {}) {
	        return new ParsedListing(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.count = source["count"];
	        this.level = source["level"];
	    }
	}

}

