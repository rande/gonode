type ResultValue = any | null;
type ResultError = Error | null;
type FuncResult = [ResultValue, ResultError];
type Arg = string | any;

let nonExistantError = new Error('Service does not exist');
let circularReferenceError = new Error('Circular reference error');

export interface ContainerInterface {
    get(name: string): FuncResult;
    set(name: string, service: any): ContainerInterface;
}

export class Definition {
    private klass: any;
    private args: Array<any>;
    private tags: string[];

    constructor(klass: any, args: Arg[] = [], tags: string[] = []) {
        this.klass = klass;
        this.args = args;
        this.tags = tags;
    }

    getArguments() {
        return this.args;
    }

    getClass() {
        return this.klass;
    }

    getTags() {
        return this.tags;
    }
}

export class ContainerBuilder {
    private definitions: Map<string, Definition>;
    private currents: Array<string>;

    constructor() {
        this.definitions = new Map<string, Definition>();
        this.currents = [];
    }

    get(name: string): Definition | undefined {
        return this.definitions.get(name);
    }

    set(name: string, definition: Definition): ContainerBuilder {
        this.definitions.set(name, definition);

        return this;
    }

    def(name: string, klass: any, args: Arg[]): ContainerBuilder {
        this.set(name, new Definition(klass, args));

        return this;
    }

    build(container: Container): FuncResult {
        let error: Error;
        this.definitions.forEach((v: Definition, k: string) => {

            let [service, err] = this.buildService(k, v);

            if (error) {
                return;
            }

            if (err) {
                error = err;
                return;
            }

            container.set(k, service);
        });

        return [container, null];
    }

    private buildService(name: string, def: Definition): FuncResult {
        let lastError = null;

        if (this.currents.includes(name)) {
            return [null, circularReferenceError];
        }

        this.currents.push(name);

        // resolve argument
        let args = def.getArguments().map((arg: any) => {
            if (typeof arg !== 'string') {
                return arg;
            }

            if (arg.length > 0 && arg[0] === '@') {
                // solve deps
                let childDefinition = this.get(arg.substring(1));

                if (!childDefinition) {
                    lastError = nonExistantError;
                    return arg;
                }

                let [child, err] = this.buildService(arg.substring(1), childDefinition);

                if (err) {
                    lastError = err;
                    return arg;
                }

                return child;
            }

            return arg;
        });

        this.currents.pop();

        let service: any;

        try {
            service = new (def.getClass())(...args);
        } catch (err) { // not a class
            service = (def.getClass())(...args);
        }

        return [service, lastError];
    }
}

export class Container implements ContainerInterface {
    private services: Map<string, any>;

    constructor() {
        this.services = new Map<string, Definition>();
    }

    get(name: string): FuncResult {
        if (name in this.services) {
            return [this.services[name], null];
        }

        return [null, new Error('Not know service')];
    }

    set(name: string, service: any): Container {
        this.services[name] = service;

        return this;
    }
}
