import { Component, OnInit, OnDestroy } from '@angular/core';
import { UiService, UiLink } from 'fgx/services/ui.service';

@Component({
  selector: 'gonode-core-breadcrumbs',
  templateUrl: './breadcrumbs.component.html',
  styleUrls: ['./breadcrumbs.component.scss']
})
export class BreadcrumbsComponent implements OnInit, OnDestroy {
  public breadcrumbs: UiLink[];
  private _breadcrumbs: any;

  constructor(private ui: UiService) {
    console.log('[BreadcrumbsComponent] construct', { ui });
    this.breadcrumbs = ui.getBreadcrumbs();
    this._breadcrumbs = ui.getBreadcrumbEmitter().subscribe(breadcrumbs => {
      console.log('[BreadcrumbsComponent] Get breadcrumbs', { breadcrumbs });

      this.breadcrumbs = breadcrumbs;
    });
  }

  ngOnInit() {}

  ngOnDestroy() {
    this._breadcrumbs.unsubscribe();
  }
}
