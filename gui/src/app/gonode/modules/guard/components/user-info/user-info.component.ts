import { Component, OnInit, OnDestroy } from '@angular/core';
import { AuthService } from 'shared/services/auth.service';
import { Credentials } from 'shared/models/credentials';

@Component({
  selector: 'gonode-guard-user-info',
  templateUrl: './user-info.component.html',
  styleUrls: ['./user-info.component.css']
})
export class UserInfoComponent implements OnInit, OnDestroy {
  credentials: Credentials;
  _auth: any;

  constructor(private auth: AuthService) {
    console.log('[UserInfoComponent] constructor');
    this.credentials = auth.credentials();
  }

  ngOnInit() {
    console.log('[UserInfoComponent] ngOnInit');

    this._auth = this.auth.current().subscribe(credentials => {
      console.log('[UserInfoComponent] New credential: ', credentials);
      this.credentials = credentials;
    });
  }

  ngOnDestroy() {
    console.log('[UserInfoComponent] ngOnDestroy ');
    this._auth.unsubscribe();
  }
}
