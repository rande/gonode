import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { KeysIteratorPipe } from './pipes/keys-iterator.pipe';
import { BreadcrumbsComponent } from './components/breadcrumbs/breadcrumbs.component';
import { UploadFieldComponent } from './components/upload-field/upload-field.component';

@NgModule({
  imports: [CommonModule, RouterModule],
  declarations: [KeysIteratorPipe, BreadcrumbsComponent, UploadFieldComponent],
  exports: [KeysIteratorPipe, BreadcrumbsComponent, UploadFieldComponent]
})
export class CoreModule {}
