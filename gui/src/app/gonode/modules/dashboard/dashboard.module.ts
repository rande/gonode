import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component';

export const DashboardRoutes: Routes = [
  { path: '',  component: DashboardComponent },
];

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [DashboardComponent]
})
export class DashboardModule { }
