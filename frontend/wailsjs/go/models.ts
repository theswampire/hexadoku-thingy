export namespace main {
	
	export class Sudoku {
	    size: number;
	    values: number[][];
	    locked: boolean[][];
	
	    static createFrom(source: any = {}) {
	        return new Sudoku(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.size = source["size"];
	        this.values = source["values"];
	        this.locked = source["locked"];
	    }
	}

}

