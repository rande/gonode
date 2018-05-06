import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';

import { HeaderComponent } from './components/header/header.component';
import { MainLayoutComponent } from './layouts/main-layout/main-layout.component';
import { LoginLayoutComponent } from './layouts/login-layout/login-layout.component';
import { MODULES } from './nebular.module';
import { TableLinkComponent } from 'fgx/components/table-cell-link/table-cell-link.component';

@NgModule({
  imports: [...MODULES, RouterModule, CommonModule],
  exports: [],
  declarations: [
    HeaderComponent,
    MainLayoutComponent,
    LoginLayoutComponent,
    TableLinkComponent
  ],
  entryComponents: [TableLinkComponent]
})
export class FgxModule {}
