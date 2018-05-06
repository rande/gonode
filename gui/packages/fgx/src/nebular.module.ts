import { NgModule } from '@angular/core';

import {
  NbActionsModule,
  NbCardModule,
  NbLayoutModule,
  NbMenuModule,
  NbRouteTabsetModule,
  NbSearchModule,
  NbSidebarModule,
  NbTabsetModule,
  NbThemeModule,
  NbUserModule,
  NbCheckboxModule,
  NbPopoverModule,
  NbContextMenuModule,
  NbProgressBarModule,
  NbSidebarService,
  NbListModule
} from '@nebular/theme';

import { NbSecurityModule } from '@nebular/security';

import { Ng2SmartTableModule as NbSmartTableModule } from 'ng2-smart-table';

export const MODULES = [
  NbActionsModule,
  NbCardModule,
  NbLayoutModule,
  NbMenuModule,
  NbRouteTabsetModule,
  NbSearchModule,
  NbSidebarModule,
  NbTabsetModule,
  NbThemeModule,
  NbUserModule,
  NbCheckboxModule,
  NbPopoverModule,
  NbContextMenuModule,
  NbProgressBarModule,
  NbSecurityModule,
  NbListModule,
  NbSmartTableModule
];

@NgModule({
  imports: MODULES,
  exports: MODULES,
  providers: [NbSidebarService]
})
export class NebularModule {}
