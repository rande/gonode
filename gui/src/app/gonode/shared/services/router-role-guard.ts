import { Injectable } from '@angular/core';
import { Router, CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from './auth.service';

@Injectable()
export class RoleGuardService implements CanActivate {

  constructor(public auth: AuthService, public router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    const expectedRole = route.data.role || 'ANONYMOUS';

    console.log('[RoleGuardService] ', expectedRole, this.auth.credentials(), state);

    if (!this.auth.credentials().hasRole(expectedRole)) {
      console.log('[RoleGuardService] invalid role, redirect to /guard/login', {route, expectedRole});

      if (state.url !== '/guard/login') {
        this.router.navigate(['/guard/login']);
      }

      return false;
    }

    console.log('[RoleGuardService] valid route', {route});

    return true;
  }
}
