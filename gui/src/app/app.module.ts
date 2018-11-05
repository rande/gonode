import { BrowserModule } from '@angular/platform-browser';
import { NgModule, Component } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { RouterModule, Routes } from '@angular/router';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NbThemeModule, NbMenuModule } from '@nebular/theme';
import { NbSecurityModule, NbRoleProvider } from '@nebular/security';
import { of as observableOf } from 'rxjs';
import { FgxModule } from 'fgx/fgx.module';

import { UserInfoComponent } from 'modules/guard/components/user-info/user-info.component';

import { RoleGuardService } from 'shared/services/router-role-guard';
import { MainLayoutComponent } from 'fgx/layouts/main-layout/main-layout.component';
import { LoginLayoutComponent } from 'fgx/layouts/login-layout/login-layout.component';

import { GuardModule, GuardRoutes } from 'modules/guard/guard.module';
import {
  DashboardRoutes,
  DashboardModule
} from 'modules/dashboard/dashboard.module';
import { MediaRoutes, MediaModule } from 'modules/media/media.module';
import { LogoutComponent } from 'modules/guard/components/logout/logout.component';
import { ErrorModule } from 'modules/error/error.module';
import { CoreModule } from 'modules/core/core.module';
import { UiService } from 'fgx/services/ui.service';
import { UploadService } from 'shared/services/upload.service';

/**
 * Main entry point with the router outlet
 *
 * <router-outlet></router-outlet>
 */
@Component({
  selector: 'app-root',
  template: '<router-outlet></router-outlet>'
})
export class AppComponent {}

export class NbSimpleRoleProvider extends NbRoleProvider {
  getRole() {
    // here you could provide any role based on any auth flow
    return observableOf('guest');
  }
}

/**
 * Configure the routes of this application
 */
const routes: Routes = [
  {
    path: 'guard/logout',
    component: LogoutComponent,
    canActivate: [RoleGuardService],
    data: { role: 'ROLE_ADMIN' }
  },
  {
    path: 'guard',
    component: LoginLayoutComponent,
    canActivate: [RoleGuardService],
    children: GuardRoutes
  },
  {
    path: 'dashboard',
    component: MainLayoutComponent,
    canActivate: [RoleGuardService],
    data: { role: 'ROLE_ADMIN' },
    children: DashboardRoutes
  },
  {
    path: 'media',
    component: MainLayoutComponent,
    canActivate: [RoleGuardService],
    data: { role: 'ROLE_ADMIN' },
    children: MediaRoutes
  },
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' }
];

@NgModule({
  exports: [RouterModule],
  imports: [RouterModule.forRoot(routes, { enableTracing: false })]
})
export class AppRoutingModule {}

/**
 * Main declaration, wrapping all hight level dependencies of this current application
 */
@NgModule({
  declarations: [AppComponent, UserInfoComponent],
  imports: [
    NbThemeModule.forRoot({ name: 'corporate' }),
    NbSecurityModule.forRoot(),
    NbMenuModule.forRoot(),
    // NebularModule,
    FgxModule,
    BrowserModule,
    GuardModule,
    DashboardModule,
    MediaModule,
    ErrorModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    CoreModule
  ],
  providers: [
    RoleGuardService,
    UiService,
    UploadService,
    {
      provide: NbRoleProvider,
      useClass: NbSimpleRoleProvider
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
