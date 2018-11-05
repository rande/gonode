import { Injectable } from '@angular/core';
import createUppy, { Uppy } from '@uppy/core';
import Dashboard from '@uppy/dashboard';
// import Tus from '@uppy/tus';
import XHRUpload from '@uppy/xhr-upload';
import { ApiService } from './api.service';

@Injectable()
export class UploadService /* extends Uppy*/ {
  readonly uppy = Uppy;

  constructor(private apiService: ApiService) {}

  configureField(node, uuid): Uppy {
    const endpoint = this.apiService.getUploadEndpoint(node);

    console.log('[UploadService:configureField]', { endpoint });

    return createUppy({
      autoProceed: false,
      debug: true,
      restrictions: {
        maxNumberOfFiles: 1,
        allowedFileTypes: ['*/*'],
        maxFileSize: 0,
        minNumberOfFiles: 0
      }
    })
      .use(Dashboard, {
        target: '.upload-' + uuid,
        replaceTargetContent: true,
        inline: true,
        proudlyDisplayPoweredByUppy: false
      })
      .use(XHRUpload, {
        endpoint: endpoint.url,
        method: endpoint.method,
        formData: false
      });
  }
}
