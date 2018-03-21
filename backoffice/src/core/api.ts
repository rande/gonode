import {stringify} from 'query-string';

// export NodeMeta: Map<string, string>;
// export NodeData: Map<string, string>;

export class Node {
    private _uuid: string = "";
    private _type: string = "";
    private _name: string = "";
    private _slug: string = "";
    private _path: string = "";
    private _status: number = 0;
    private _weight: number = 0;
    private _revision: number = 0;
    private _createdAt: string  = "";
    private _updatedAt: string  = "";
    private _enabled: boolean = false;
    private _deleted: boolean = false;
    private _parents: Array<string> = [];
    private _updatedBy: string = "";
    private _createdBy: string = "";
    private _parentUuid: string = "";
    private _source: string = "";
    private _modules: Array<string> = [];
    private _access: Array<string> = [];
    private _meta: Map<string, any> = new Map<string, any>();
    private _data: Map<string, any> = new Map<string, any>();

    get uuid(): string {
        return this._uuid;
    }

    set uuid(value: string) {
        this._uuid = value;
    }

    get type(): string {
        return this._type;
    }

    set type(value: string) {
        this._type = value;
    }

    get name(): string {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    get slug(): string {
        return this._slug;
    }

    set slug(value: string) {
        this._slug = value;
    }

    get path(): string {
        return this._path;
    }

    set path(value: string) {
        this._path = value;
    }

    get status(): number {
        return this._status;
    }

    set status(value: number) {
        this._status = value;
    }

    get weight(): number {
        return this._weight;
    }

    set weight(value: number) {
        this._weight = value;
    }

    get revision(): number {
        return this._revision;
    }

    set revision(value: number) {
        this._revision = value;
    }

    get createdAt(): string {
        return this._createdAt;
    }

    set createdAt(value: string) {
        this._createdAt = value;
    }

    get updatedAt(): string {
        return this._updatedAt;
    }

    set updatedAt(value: string) {
        this._updatedAt = value;
    }

    get enabled(): boolean {
        return this._enabled;
    }

    set enabled(value: boolean) {
        this._enabled = value;
    }

    get deleted(): boolean {
        return this._deleted;
    }

    set deleted(value: boolean) {
        this._deleted = value;
    }

    get parents(): Array<string> {
        return this._parents;
    }

    set parents(value: Array<string>) {
        this._parents = value;
    }

    get updatedBy(): string {
        return this._updatedBy;
    }

    set updatedBy(value: string) {
        this._updatedBy = value;
    }

    get createdBy(): string {
        return this._createdBy;
    }

    set createdBy(value: string) {
        this._createdBy = value;
    }

    get parentUuid(): string {
        return this._parentUuid;
    }

    set parentUuid(value: string) {
        this._parentUuid = value;
    }

    get source(): string {
        return this._source;
    }

    set source(value: string) {
        this._source = value;
    }

    get modules(): Array<string> {
        return this._modules;
    }

    set modules(value: Array<string>) {
        this._modules = value;
    }

    get access(): Array<string> {
        return this._access;
    }

    set access(value: Array<string>) {
        this._access = value;
    }

    get meta(): Map<string, any> {
        return this._meta;
    }

    set meta(value: Map<string, any>) {
        this._meta = value;
    }

    get data(): Map<string, any> {
        return this._data;
    }

    set data(value: Map<string, any>) {
        this._data = value;
    }
}

export class Query {
    private _page: number;
    private _perPage: number;
    private _orderBy: Array<string>;
    private _types: Array<string>;
    private _uuid: Array<string>;
    private _status: Array<string>;
    private _weight: Array<string>;
    private _enabled: boolean;
    private _deleted: boolean;
    private _updatedBy: Array<string>;
    private _createdBy: Array<string>;
    private _parentUuid: Array<string>;
    private _source: Array<string>;
    private _meta: Map<string, Array<string>>;
    private _data: Map<string, Array<string>>;

    constructor(page = 1, perPage = 32) {
        this._page = page;
        this._perPage = perPage;
        this._orderBy = [];
        this._types = [];
        this._uuid = [];
        this._status = [];
        this._weight = [];
        this._updatedBy = [];
        this._createdBy = [];
        this._parentUuid = [];
        this._source = [];
        this._meta = new Map();
        this._data = new Map();
        this._enabled = true;
        this._deleted = false;
    }

    setPage(page: number): Query {
        this._page = page;

        return this;
    }

    setPerPage(perPage: number): Query {
        this._perPage = perPage;

        return this;
    }

    addOrderBy(field: string, order: string = 'ASC'): Query {
        this._orderBy.push(`${field},${order}`);

        return this;
    }

    addType(type: string): Query {
        this._types.push(type);

        return this;
    }

    addUuid(uuid: string): Query {
        this._uuid.push(uuid);

        return this;
    }

    setUuid(uuid: Array<string>): Query {
        this._uuid = uuid;

        return this;
    }

    setEnabled(enabled: boolean): Query {
        this._enabled = enabled;

        return this;
    }

    setDeleted(deleted: boolean): Query {
        this._deleted = deleted;

        return this;
    }

    addUpdatedBy(updatedBy: string): Query {
        this._updatedBy.push(updatedBy);

        return this;
    }

    addCreatedBy(createdBy: string): Query {
        this._createdBy.push(createdBy);

        return this;
    }

    addData(field: string, value: any): Query {
        if (!this._data.has(field)) {
            this._data.set(field, []);
        }

        let data = this._data.get(field);
        if (data) {
            data.push(value)
        }

        return this
    }

    addMeta(field: string, value: any): Query {
        if (!this._meta.has(field)) {
            this._meta.set(field, []);
        }

        let meta = this._data.get(field);
        if (meta) {
            meta.push(value)
        }

        return this
    }

    get page(): number {
        return this._page;
    }

    get perPage(): number {
        return this._perPage;
    }

    get orderBy(): Array<string> {
        return this._orderBy;
    }

    get types(): Array<string> {
        return this._types;
    }

    get uuid(): Array<string> {
        return this._uuid;
    }

    get status(): Array<string> {
        return this._status;
    }

    get weight(): Array<string> {
        return this._weight;
    }

    get enabled(): boolean {
        return this._enabled;
    }

    get deleted(): boolean {
        return this._deleted;
    }

    get updatedBy(): Array<string> {
        return this._updatedBy;
    }

    get createdBy(): Array<string> {
        return this._createdBy;
    }

    get parentUuid(): Array<string> {
        return this._parentUuid;
    }

    get source(): Array<string> {
        return this._source;
    }

    get meta(): Map<string, Array<string>> {
        return this._meta;
    }

    get data(): Map<string, Array<string>> {
        return this._data;
    }
}

export interface Pager {
    readonly page: number;
    readonly perPage: number;
    readonly items: Array<Node>;
    readonly previous: number;
    readonly next: number;
}

export class ApiClient {
    private readonly _url: string;
    private readonly _version: string;
    private _token: string;

    /**
     *
     * @param {string} url
     * @param {string} version
     */
    constructor(url: string, version: string = "1.0") {
        this._url = url;
        this._version = version;
        this._token = "";
    }

    /**
     *
     * @param {Query} query
     * @returns {Pager}
     */
    async find(query: Query): Promise<Pager> {
        let qs = {
            per_page: query.perPage,
            page: query.page,
            type: query.types,
            uuid: query.uuid,
            status: query.status,
            weight: query.weight,
            enabled: query.enabled,
            deleted: query.deleted,
            updated_by: query.updatedBy,
            created_by: query.createdBy,
            parent_uuid: query.parentUuid,
            source: query.source,
            order_by: query.orderBy,
        };

        query.meta.forEach((v, k) => {
            qs[`meta.${k}`] = v;
        });

        query.data.forEach((v, k) => {
            qs[`data.${k}`] = v;
        });

        let response = await this.doAuthRequest('GET', `/nodes?${stringify(qs)}`);

        return {
            page: response.page,
            perPage: response.per_page,
            items: response.elements,
            previous: response.previous,
            next: response.next,
        };
    }

    private async doAuthRequest(method: string, url: string, body = null): Promise<any> {
        
        let response: Response;

        try {
            response = await fetch(`${this._url}/v${this._version}/${url}`, {
                method: method,
                headers: {
                    Authorization: `Bearer ${this._token}`
                }
            });    
        } catch(e) {
            return false;
        }

        if (!response.ok) {
            return false;
        }

        return await response.json();
    }

    async findOne(query: Query): Promise<Node|boolean> {
        let pager = await this.find(query);

        if (pager.items.length != 1) {
            return false;
        }

        return pager.items[0];
    }

    /**
     *
     * @param {string} username
     * @param {string} password
     * @returns {Promise<boolean>}
     */
    async signin(username: string, password: string): Promise<boolean> {
        const params = new URLSearchParams();
        params.append('username', username);
        params.append('password', password);

        let response: Response;

        try {
            response = await fetch(`${this._url}/v${this._version}/login`, {
                method: 'POST',
                body: params
            });
        } catch (e) {
            console.error(e);
            
            return false;
        }

        if (!response.ok) {
            return false;
        }

        const auth = await (response.json() as Promise<AuthResult>);

        if (auth.status === "OK") {
            this._token = auth.token;

            return true;
        }

        return false;
    }

    /**
     *
     * @returns {boolean}
     */
    isAuthenticated() {
        return this._token.length > 0;
    }

    /**
     *
     * @returns {Promise<boolean>}
     */
    async logout(): Promise<boolean> {
        this._token = "";

        return true;
    }
}

export interface AuthResult {
    message: string,
    status: string,
    token: string
}
