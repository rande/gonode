import { Component, OnChanges, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UnableToAuthenticateError } from 'fgx/types/errors';
import { AuthService } from 'shared/services/auth.service';

@Component({
  selector: 'gonode-guard-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit, OnChanges {
  loginForm: FormGroup;
  isLoading = false;
  globalErrorMessage = '';

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {
    this.createForm();
  }

  ngOnInit() {
    console.log('[LoginComponent] ngOnInit');
  }

  createForm() {
    this.loginForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required]
    });

    this.rebuildForm();
  }

  rebuildForm() {
    this.loginForm.reset({
      username: 'admin',
      password: 'admin'
    });
  }

  ngOnChanges() {
    console.log('[LoginComponent] ngOnChanges');

    this.rebuildForm();
  }

  async onSubmit() {
    const formModel = this.loginForm.value;

    // result values
    this.isLoading = true;
    this.globalErrorMessage = '';

    const result = await this.auth.login(
      formModel.username,
      formModel.password
    );

    this.isLoading = false;

    if (result.error) {
      this.globalErrorMessage = 'An unexpected error occurs on the server';
    }

    if (result.error === UnableToAuthenticateError) {
      this.globalErrorMessage =
        'Unable to authenticate the user, please check username and password';
    }

    if (!result.value) {
      console.log('[LoginComponent] Authentification failed');

      return;
    }

    console.log(
      '[LoginComponent] Authentification successed, navigate to /dashboard'
    );
    await this.router.navigate(['/dashboard']);
  }

  revert() {
    this.rebuildForm();
  }
}
