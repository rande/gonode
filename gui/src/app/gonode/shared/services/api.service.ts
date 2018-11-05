import { Injectable, EventEmitter } from '@angular/core';
import {
  HttpClient,
  HttpParams,
  HttpHeaders,
  HttpResponse,
  HttpErrorResponse
} from '@angular/common/http';

import {
  NoResultFoundError,
  InvalidStatusCodeError,
  UnableToAuthenticateError,
  UnexpectedHttpError,
  UnauthorizedHttpError,
  ForbiddenHttpError,
  MutationError,
  UnprocessableEntityError,
  VersionConflictError
} from 'fgx/types/errors';

import { Maybe, Result } from 'fgx/types';

export const NodeStatus = [
  [0, 'New'],
  [1, 'Draft'],
  [2, 'Completed'],
  [3, 'Validated']
];

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

    const data = this._data.get(field);
    if (data) {
      data.push(value);
    }

    return this;
  }

  addMeta(field: string, value: any): Query {
    if (!this._meta.has(field)) {
      this._meta.set(field, []);
    }

    const meta = this._data.get(field);
    if (meta) {
      meta.push(value);
    }

    return this;
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
  readonly nodes: Array<Node>;
  readonly previous: number;
  readonly next: number;
}

interface ApiPager {
  readonly page: number;
  readonly per_page: number;
  readonly elements: Array<Node>;
  readonly previous: number;
  readonly next: number;
}

export interface Node {
  readonly uuid: string;
  readonly type: string;
  readonly name: string;
  readonly slug: string;
  readonly path: string;
  readonly status: number;
  readonly weight: number;
  readonly revision: number;
  readonly createdAt: Date;
  readonly updatedAt: Date;
  readonly enabled: boolean;
  readonly deleted: boolean;
  readonly parents: Array<string>;
  readonly updatedBy: string;
  readonly createdBy: string;
  readonly parentUuid: string;
  readonly source: string;
  readonly modules: Array<string>;
  readonly access: Array<string>;
  readonly meta: {
    [key: string]: any;
  };
  readonly data: {
    [key: string]: any;
  };
}

export function normalizeNode(data: object): Node {
  return {
    uuid: '',
    type: '',
    name: '',
    slug: '',
    path: '',
    status: 0,
    weight: 0,
    revision: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
    enabled: false,
    deleted: false,
    parents: [],
    updatedBy: '',
    createdBy: '',
    parentUuid: '',
    source: '',
    modules: [],
    access: [],
    meta: {},
    data: {},
    ...data
  };
}

function normalizePager(pager: ApiPager): Pager {
  return {
    page: pager.page,
    perPage: pager.per_page,
    previous: pager.previous,
    next: pager.next,
    nodes: pager.elements.map(item => normalizeNode(item))
  };
}

export function createPager(data = {}) {
  return {
    page: 1,
    perPage: 0,
    nodes: [],
    previous: 1,
    next: 0
  };
}

export interface AuthCredentials {
  roles: string[];
  username: string;
}

export interface AuthResult {
  message: string;
  status: string;
  credentials: AuthCredentials;
}

export function createNode(type: string, data = {}): Node {
  return {
    uuid: '',
    type,
    name: '',
    slug: '',
    path: '',
    status: 0,
    weight: 1,
    revision: 1,
    createdAt: new Date(),
    updatedAt: new Date(),
    enabled: false,
    deleted: false,
    parents: [],
    updatedBy: '',
    createdBy: '',
    parentUuid: '',
    source: '',
    modules: [],
    access: [],
    meta: {},
    data: {},
    ...data
  };
}

export function updateNode(node: Node, data: {}): Maybe<Node> {
  if ('type' in data) {
    return Result(node, MutationError);
  }

  if ('uuid' in data) {
    return Result(node, MutationError);
  }

  return Result({
    ...node,
    ...data
  });
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private url = 'http://localhost:2508/api/v1.0';

  constructor(private http: HttpClient) {}

  async login(
    username: string,
    password: string
  ): Promise<Maybe<AuthCredentials>> {
    const req = this.http.post<AuthResult>(
      `${this.url}/login`,
      { username, password },
      {
        headers: new HttpHeaders({
          'Content-Type': 'application/json'
        }),
        observe: 'response'
      }
    );

    let response: HttpResponse<AuthResult>;

    try {
      response = await req.toPromise();
    } catch (err) {
      console.log('[ApiService] unable to log', { err });

      return Result(null, UnableToAuthenticateError, err);
    }

    if (!response.ok) {
      return Result(null, InvalidStatusCodeError);
    }

    if (response.body.status === 'OK') {
      return Result(response.body.credentials);
    }

    return Result(null, UnableToAuthenticateError);
  }

  async logout(): Promise<Maybe<boolean>> {
    // todo: send a logout request
    return Result(true);
  }

  async find(query: Query): Promise<Maybe<Pager>> {
    let params = new HttpParams()
      .set('per_page', query.perPage.toString())
      .set('page', query.page.toString())
      .set('enabled', query.enabled.toString())
      .set('deleted', query.deleted.toString())
      .set('updated_by', query.updatedBy.toString())
      .set('created_by', query.createdBy.toString());

    query.types.forEach(s => {
      params = params.append('type', s);
    });
    query.uuid.forEach(s => {
      params = params.append('uuid', s);
    });
    query.status.forEach(s => {
      params = params.append('status', s);
    });
    query.weight.forEach(s => {
      params = params.append('weight', s);
    });
    query.parentUuid.forEach(s => {
      params = params.append('parent_uuid', s);
    });
    query.source.forEach(s => {
      params = params.append('source', s);
    });
    query.orderBy.forEach(s => {
      params = params.append('order_by', s);
    });

    query.meta.forEach((v, k) => {
      v.forEach(s => {
        params = params.append(`meta.${k}`, s);
      });
    });

    query.data.forEach((v, k) => {
      v.forEach(s => {
        params = params.append(`data.${k}`, s);
      });
    });

    const result = await this.doAuthRequest<ApiPager>('GET', '/nodes', params);

    if (result.error) {
      return Result<Pager>(undefined, result.error, result.originalError);
    }

    return Result<Pager>(normalizePager(result.value));
  }

  async findUuid(uuid: string) {
    return await this.doAuthRequest<Node>('GET', `/nodes/${uuid}`);
  }

  async findOne(query: Query) {
    const result = await this.find(query);

    if (result.error) {
      return Result<Node>(undefined, result.error, result.originalError);
    }

    if (result.value.nodes.length !== 1) {
      return Result<Node>(undefined, NoResultFoundError, result.originalError);
    }

    return Result<Node>(result.value.nodes[0]);
  }

  async save(node: Node) {
    let method = 'PUT';
    let url = `/nodes/${node.uuid}`;

    if (node.uuid.length === 0) {
      method = 'POST';
      url = '/nodes';
    }

    return this.doAuthRequest<Node>(method, url, undefined, node);
  }

  private async doAuthRequest<T>(
    method: string,
    url: string,
    params = new HttpParams(),
    body?: any
  ): Promise<Maybe<T>> {
    // to remove on production
    // params = params.set('access_token', this._token);

    console.log(method, url, params, body, url);

    // debugger;
    const req = this.http.request<T>(method, `${this.url}${url}`, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      }),
      observe: 'response',
      params: params,
      body: body
    });

    let response: HttpResponse<T>;

    try {
      response = await req.toPromise();
    } catch (err) {
      if (err instanceof HttpErrorResponse) {
        switch (err.status) {
          case 401:
            return Result<T>(undefined, UnauthorizedHttpError, err);
          case 403:
            return Result<T>(undefined, ForbiddenHttpError, err);
          case 409:
            return Result<T>(undefined, VersionConflictError, err);
          case 422:
            return Result<T>(undefined, UnprocessableEntityError, err);
        }
      }

      console.log(
        '[ApiService:doAuthRequest] unexpected http error to perform HTTP Request',
        {
          err
        }
      );

      return Result<T>(undefined, UnexpectedHttpError, err);
    }

    if (!response.ok) {
      return Result<T>(undefined, InvalidStatusCodeError);
    }

    return Result<T>(response.body);
  }

  getUploadEndpoint(node: Node) {
    return {
      url: `${this.url}/nodes/${node.uuid}?raw`,
      method: 'PUT'
    };
  }
}
