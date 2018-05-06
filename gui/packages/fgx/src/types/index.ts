import { HttpErrorResponse } from '@angular/common/http';

export interface Pager<T> {
  readonly page: number;
  readonly perPage: number;
  readonly elements: Array<T>;
  readonly previous: number;
  readonly next: number;
}

export class Maybe<T> {
  constructor(
    public readonly value: T | undefined,
    public readonly error: Error | undefined,
    public readonly originalError: Error | undefined
  ) {}
}

// tslint:disable-next-line:no-shadowed-variable
export function Result<T>(
  value: T | undefined,
  error?: Error,
  originalError?: Error
): Maybe<T> {
  return new Maybe<T>(value, error, originalError);
}

export class HttpUnprocessableEntityError extends HttpErrorResponse {
  error: {
    [key: string]: Array<string>;
  };
  ok: boolean;
}
