import { Injectable, EventEmitter, Inject } from '@angular/core';
import { Credentials, getAnonymousCredentials } from '../models/credentials';
import { ApiService } from './api.service';
import { StorageService } from './storage.service';
import { Maybe, Result } from 'fgx/types';
import { UnableToAuthenticateError } from 'fgx/types/errors';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private _credentials = getAnonymousCredentials();
  private _credentialsEmitter = new EventEmitter<Credentials>();

  constructor(private api: ApiService, private storage: StorageService) {
    console.log('[AuthService] constructor');

    this._credentials = this.getCredentials(
      this.storage.get('gonode.user.credentials', false)
    );

    this._credentialsEmitter.emit(this._credentials);
  }

  private getCredentials(data) {
    if (!data) {
      console.log('[AuthService] emit getAnonymousCredentials()');

      return getAnonymousCredentials();
    }

    return new Credentials(data.username, data.roles);
  }

  credentials() {
    return this._credentials;
  }

  current() {
    return this._credentialsEmitter;
  }

  async logout() {
    this._credentials = getAnonymousCredentials();

    return await this.api.logout();
  }

  async login(username: string, password: string): Promise<Maybe<Credentials>> {
    const result = await this.api.login(username, password);

    // unable to login
    if (result.error) {
      return Result(null, UnableToAuthenticateError, result.error);
    }

    this._credentials = new Credentials(
      result.value.username,
      result.value.roles
    );

    this.storage.set('gonode.user.credentials', {
      username: result.value.username,
      roles: result.value.roles
    });

    this._credentialsEmitter.emit(this._credentials);

    return Result(this._credentials);
  }
}
