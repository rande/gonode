import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { ErrorModule } from 'modules/error/error.module';
import { CoreModule } from 'modules/core/core.module';

import { ListComponent } from './components/list/list.component';
import { EditComponent } from './components/edit/edit.component';
import { NebularModule } from 'fgx/nebular.module';

export const MediaRoutes: Routes = [
  { path: 'edit/:uuid', component: EditComponent },
  { path: '', component: ListComponent }
];

@NgModule({
  imports: [
    ErrorModule,
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    CoreModule,
    NebularModule
  ],
  declarations: [ListComponent, EditComponent]
})
export class MediaModule {}
