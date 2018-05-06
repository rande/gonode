import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Http404Component } from './components/http404/http404.component';
import { DispatcherComponent } from './components/dispatcher/dispatcher.component';
import { InputErrorComponent } from './components/input-error/input-error.component';
import { CoreModule } from '../core/core.module';

@NgModule({
  imports: [CommonModule, CoreModule],
  exports: [Http404Component, DispatcherComponent, InputErrorComponent],
  declarations: [Http404Component, DispatcherComponent, InputErrorComponent]
})
export class ErrorModule {}
