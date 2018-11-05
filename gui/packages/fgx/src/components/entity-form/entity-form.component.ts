import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

import {
  NoResultFoundError,
  UnprocessableEntityError
} from '../../types/errors';

import { HttpUnprocessableEntityError, Maybe, Result } from '../../types';

import { UiService } from '../../services/ui.service';

export abstract class EntityFormComponent<T> implements OnInit {
  public isLoading: boolean;
  public entity: T;
  public form: FormGroup;
  public globalErrorMessages: string[];

  private fb: FormBuilder;
  private ui: UiService;

  constructor(fb: FormBuilder, ui: UiService) {
    this.fb = fb;
    this.ui = ui;
  }

  protected abstract loadEmptyEntity();

  protected abstract async loadEntity();

  protected abstract transformFormToEntity(
    form: FormGroup,
    entity: T
  ): Maybe<T>;

  protected abstract async saveEntity(entity: T);

  protected async handleLoadError(result: Maybe<T>): Promise<boolean> {
    if (result.error === NoResultFoundError) {
      this.globalErrorMessages = ['Unable to find the node element'];

      return true;
    }

    if (result.error) {
      this.globalErrorMessages = [
        `An unexpected error occurs on the server: ${
          result.originalError ? result.originalError : result.error
        }`
      ];

      return true;
    }

    return false;
  }

  protected async handleSaveError(result: Maybe<T>): Promise<boolean> {
    if (result.error === UnprocessableEntityError) {
      const remoteErrors = result.originalError as HttpUnprocessableEntityError;

      Object.keys(remoteErrors.error).forEach(k => {
        const field = this.form.get(k);

        if (field) {
          let errors = field.errors;

          remoteErrors.error[k].forEach(e => {
            const msg = `Server-side validation: ${e}`;
            errors = { ...errors, [msg]: msg };
          });

          field.setErrors(errors);
        } else {
          remoteErrors.error[k].forEach(e => {
            this.globalErrorMessages.push(e);
          });
        }
      });

      return true;
    }

    if (result.error) {
      this.globalErrorMessages = [
        `An unexpected error occurs on the server: ${
          result.originalError ? result.originalError : result.error
        }`
      ];

      return true;
    }

    return false;
  }

  async ngOnInit() {
    this.globalErrorMessages = [];

    this.entity = this.loadEmptyEntity();

    console.log('[EditComponent:ngOnInit] init default object', {
      entity: this.entity
    });

    this.createForm();

    this.isLoading = true;

    this.configureBreadcrumbs();

    const result = await this.loadEntity();

    this.isLoading = false;

    if ((await this.handleLoadError(result)) === true) {
      // error has been handled, return
      return;
    }

    this.entity = result.value;

    console.log('[EditComponent:ngOnInit] Entity has been loaded', {
      entity: this.entity
    });

    this.configureBreadcrumbs();
    this.rebuildForm();
  }

  configureBreadcrumbs() {
    const breadcrumbs = [];

    this.ui.setBreadcrumbs(breadcrumbs);
  }

  createForm() {
    this.form = this.fb.group({});
  }

  rebuildForm() {
    this.form.reset({});
  }

  async onSubmit() {
    this.globalErrorMessages = [];

    let result = await this.transformFormToEntity(this.form, this.entity);

    if (result.error) {
      this.globalErrorMessages = ['Unable to locally update node'];

      return;
    }

    // result values
    this.isLoading = true;

    result = await this.saveEntity(result.value);

    this.isLoading = false;

    if ((await this.handleSaveError(result)) === true) {
      // error has been handled, return
      return;
    }

    this.entity = result.value;
    this.rebuildForm();
  }
}
