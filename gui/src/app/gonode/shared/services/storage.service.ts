import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class StorageService {
  prefix: string;

  constructor() {
    this.prefix = 'gonode/';
  }

  get<T>(key: string, dflt?: T): T | undefined {
    const value = localStorage.getItem(`${this.prefix}${key}`);

    if (!value) {
      return dflt;
    }

    return JSON.parse(value);
  }

  set(key: string, value: any): void {
    localStorage.setItem(`${this.prefix}${key}`, JSON.stringify(value));
  }
}
