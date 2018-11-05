import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

import {
  VersionConflictError,
  UnauthorizedHttpError,
  ForbiddenHttpError
} from 'fgx/types/errors';

import {
  ApiService,
  Node,
  createNode,
  updateNode,
  normalizeNode
} from 'shared/services/api.service';
import { nodeStatusValidator } from 'shared/validators';
import { UiService } from 'fgx/services/ui.service';
import { AuthService } from 'shared/services/auth.service';
import { EntityFormComponent } from 'fgx/components/entity-form/entity-form.component';
import { Result, Maybe } from 'fgx/types';

@Component({
  selector: 'gonode-media-edit',
  templateUrl: './edit.component.html',
  styleUrls: ['./edit.component.css']
})
export class EditComponent extends EntityFormComponent<Node> {
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private api: ApiService,
    private auth: AuthService,
    private ui: UiService,
    private fb: FormBuilder
  ) {
    super(fb, ui);
  }

  protected loadEmptyEntity() {
    return createNode('media.image');
  }

  protected async loadEntity() {
    return await this.api.findUuid(this.route.snapshot.paramMap.get('uuid'));
  }

  protected async saveEntity(entity) {
    return await this.api.save(entity);
  }

  protected async handleLoadError(result): Promise<boolean> {
    if (
      result.error === UnauthorizedHttpError ||
      result.error === ForbiddenHttpError
    ) {
      await this.auth.logout();

      return this.router.navigateByUrl('/guard/login');
    }

    return super.handleLoadError(result);
  }

  protected async handleSaveError(result): Promise<boolean> {
    if (
      result.error === UnauthorizedHttpError ||
      result.error === ForbiddenHttpError
    ) {
      await this.auth.logout();

      return this.router.navigateByUrl('/guard/login');
    }

    if (result.error === VersionConflictError) {
      this.globalErrorMessages = [
        'A new version is available on server, please reload the form.'
      ];

      return true;
    }

    return super.handleSaveError(result);
  }

  protected transformFormToEntity(form: FormGroup, entity: Node): Maybe<Node> {
    return Result(
      normalizeNode({
        ...entity,
        ...form.value,
        data: {
          ...entity.data,
          ...(form.value.data ? form.value.data : {})
        },
        meta: {
          ...entity.meta,
          ...(form.value.meta ? form.value.meta : {})
        }
      })
    );
  }

  configureBreadcrumbs() {
    const breadcrumbs = [];

    breadcrumbs.push({
      link: '/media',
      title: 'Media'
    });

    if (this.isLoading) {
      breadcrumbs.push({
        link: '/media',
        title: 'Loading ...'
      });
    } else if (this.entity) {
      breadcrumbs.push({
        link: `/media/edit/${this.entity.uuid}`,
        title: this.entity.name
      });
    }

    this.ui.setBreadcrumbs(breadcrumbs);
  }

  createForm() {
    this.form = this.fb.group({
      name: ['', [Validators.required]],
      slug: ['', [Validators.required]],
      status: ['', [Validators.required, nodeStatusValidator()]],
      weight: [
        '',
        [Validators.required, Validators.min(0), Validators.max(100)]
      ],
      enabled: ['', [Validators.required]],
      data: this.fb.group({
        reference: ['', [Validators.required]],
        name: ['', [Validators.required]],
        source_url: ['']
      })
    });
  }

  rebuildForm() {
    this.form.reset({
      name: this.entity.name,
      slug: this.entity.slug,
      status: this.entity.status,
      weight: this.entity.weight,
      enabled: this.entity.enabled,
      data: {
        reference: this.entity.data.reference,
        name: this.entity.data.name,
        source_url: this.entity.data.source_url
      }
    });
  }
}
