import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { Routes } from '@angular/router';

import { LoginComponent } from './components/login/login.component';
import { LogoutComponent } from './components/logout/logout.component';
import { LostPasswordComponent } from './components/lost-password/lost-password.component';

import { NebularModule } from 'fgx/nebular.module';

export const GuardRoutes: Routes = [
  { path: 'login', component: LoginComponent },
  // { path: 'logout',  component: LogoutComponent }, => Cannot be there as the RoleGuard will loop
  { path: 'lost-password', component: LostPasswordComponent }
];

@NgModule({
  imports: [CommonModule, FormsModule, ReactiveFormsModule, NebularModule],
  declarations: [LoginComponent, LostPasswordComponent, LogoutComponent]
})
export class GuardModule {}
