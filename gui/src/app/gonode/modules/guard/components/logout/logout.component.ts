import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService } from 'shared/services/auth.service';

@Component({
  template: `<p>logout works!</p>`
})
export class LogoutComponent implements OnInit {
  constructor(private auth: AuthService, private router: Router) {
    console.log('[LogoutComponent] constructor');
  }

  async ngOnInit() {
    console.log('[LogoutComponent] ngOnInit');

    if (await this.auth.logout()) {
      console.log('[LogoutComponent] logout successful');
      this.router.navigate(['/guard/login']);
    }
  }
}
