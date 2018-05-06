import { Injectable, EventEmitter } from '@angular/core';

export interface UiLink {
  title: string;
  link: string;
}

@Injectable({
  providedIn: 'root'
})
export class UiService {
  private title: string;
  private breadcrumbs: UiLink[];

  private _breadcrumbsEmitter = new EventEmitter();

  constructor() {
    this.title = 'Admin';
    this.setBreadcrumbs([
      {
        title: 'Dashboard',
        link: '/dashboard'
      }
    ]);
  }

  public setBreadcrumbs(breadcrumbs: UiLink[]) {
    this.breadcrumbs = [
      {
        title: '🏠',
        link: '/'
      },
      ...breadcrumbs
    ];

    console.log('[UiService] setBreadcrumbs', { breadcrumbs });

    this._breadcrumbsEmitter.emit(this.breadcrumbs);
  }

  public getBreadcrumbEmitter() {
    return this._breadcrumbsEmitter;
  }

  public getBreadcrumbs() {
    return this.breadcrumbs;
  }
}
