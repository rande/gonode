import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { UnauthorizedHttpError, ForbiddenHttpError } from 'fgx/types/errors';

import { ApiService, Query, Pager } from 'shared/services/api.service';
import { UiService } from 'fgx/services/ui.service';
import { TableLinkComponent } from 'fgx/components/table-cell-link/table-cell-link.component';
import { DataSource } from 'ng2-smart-table/lib/data-source/data-source';
import { AuthService } from 'shared/services/auth.service';

export interface RemotePager {
  readonly page: number;
  readonly perPage: number;
  readonly data: Array<any>;
  readonly previous: number;
  readonly next: number;
}

export interface DataSourceFilters {
  page: number;
  perPage: number;
  filters: object;
}

class RemoteDataSource extends DataSource {
  protected pager: RemotePager;
  protected filters: DataSourceFilters;
  protected loadCallback: Function;

  constructor(loadCallback: Function, filters: DataSourceFilters) {
    super();

    this.loadCallback = loadCallback;
    this.filters = filters;
  }

  getAll(): Promise<any> {
    // console.log('RemoteDataSource:getAll');
    return this.getElements();
  }

  async getElements(): Promise<any> {
    // console.log('RemoteDataSource:getElements');

    this.pager = await this.loadCallback(this.filters);

    return Promise.resolve(this.pager.data.slice(0));
  }

  getSort() {
    // console.log('RemoteDataSource:getSort');
    return [];
  }

  getFilter() {
    // console.log('RemoteDataSource:getFilter');
    return {
      filters: [],
      andOperator: true
    };
  }

  getPaging() {
    // console.log('RemoteDataSource:getPaging');
    return {
      perPage: this.pager.perPage,
      page: this.pager.page
    };
  }

  count(): number {
    if (!this.pager) {
      // console.log('RemoteDataSource:count pager not initialized');
      return 0;
    }

    const count =
      this.pager.page * this.pager.perPage + (this.pager.next > 0 ? 1 : 0);

    // console.log('RemoteDataSource:count', { count });

    return count;
  }

  refresh() {
    // console.log('RemoteDataSource:refresh');
    this.emitOnChanged('refresh');
  }

  setPage(page: number, doEmit?: boolean) {
    // console.log('RemoteDataSource:setPage', { page, doEmit });
    this.filters.page = page;

    if (doEmit) {
      this.emitOnChanged('page');
    }

    this.refresh();
  }
}

@Component({
  selector: 'gonode-media-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ListComponent implements OnInit {
  source: DataSource;
  loading: boolean;
  displayedColumns = ['uuid', 'name', 'weight'];
  pager?: Pager;

  tableSettings = {
    columns: {
      uuid: {
        title: 'Uuid',
        type: 'custom',
        editable: false,
        filter: false,
        renderComponent: TableLinkComponent
      },
      type: {
        title: 'Type',
        editable: false,
        filter: false
      },
      status: {
        title: 'Status',
        editable: false,
        filter: false
      }
    },
    actions: {
      add: false,
      edit: false,
      delete: false,
      columns: false
    }
  };

  constructor(
    private api: ApiService,
    private ui: UiService,
    private auth: AuthService,
    private router: Router
  ) {}

  async ngOnInit() {
    this.loading = true;

    /**
     * Bridge the API with the RemoteDataSource
     *
     * @param filters
     */
    const loadCallback = async (filters: DataSourceFilters) => {
      // create the query.
      const query = new Query(filters.page, filters.perPage).addType(
        'media.image'
      );

      // call the api and get back the result.
      const result = await this.api.find(query);

      if (
        result.error === UnauthorizedHttpError ||
        result.error === ForbiddenHttpError
      ) {
        await this.auth.logout();

        return this.router.navigateByUrl('/guard/login');
      }

      if (result.error) {
        return {
          perPage: filters.perPage,
          page: filters.page,
          data: [],
          next: 0,
          previous: 0
        };
      }

      // return a pager with filtered data ready to be used in the table.
      return {
        perPage: result.value.perPage,
        next: result.value.next,
        previous: result.value.previous,
        page: result.value.page,
        data: result.value.nodes.map(node => ({
          uuid: { href: `/media/edit/${node.uuid}`, text: node.name },
          type: node.type,
          status: node.status
        }))
      };
    };

    this.source = new RemoteDataSource(loadCallback, {
      perPage: 32,
      page: 1,
      filters: {}
    });

    this.ui.setBreadcrumbs([
      {
        link: '/media',
        title: 'Media'
      }
    ]);

    this.loading = false;
  }
}
