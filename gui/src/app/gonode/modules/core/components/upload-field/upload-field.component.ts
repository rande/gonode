import {
  Component,
  OnInit,
  AfterViewInit,
  AfterViewChecked,
  ViewEncapsulation,
  Input,
  ChangeDetectionStrategy,
  OnChanges,
  SimpleChanges
} from '@angular/core';
import { Subject } from 'rxjs';
import { UploadService } from 'shared/services/upload.service';
import { v4 } from 'uuid';
import { Node } from 'shared/services/api.service';
import { Uppy } from '@uppy/core';

export type UploadPluginConfigurations = [String, any][];

const uppyEvents = [
  'file-added',
  'file-removed',
  'upload',
  'upload-progress',
  'upload-success',
  'complete',
  'upload-error',
  'info-visible',
  'info-hidden'
];

@Component({
  selector: 'gonode-upload-field',
  templateUrl: './upload-field.component.html',
  styleUrls: ['./upload-field.component.scss'],
  encapsulation: ViewEncapsulation.None
  // changeDetection: ChangeDetectionStrategy.Default
})
export class UploadFieldComponent
  implements OnInit, AfterViewChecked, AfterViewInit, OnChanges {
  @Input()
  node: Node;
  @Input()
  on: Subject<[string, any, any, any]>;

  uuid = v4();
  uppyInstance: Uppy;

  constructor(private uploadService: UploadService) {}

  ngOnInit() {
    // console.log('[UploadFieldComponent:ngOnInit]');
  }

  ngOnChanges(changes: SimpleChanges) {
    console.log('[UploadFieldComponent:ngOnChanges]', { changes });

    if (
      changes.node.firstChange ||
      changes.node.currentValue.uuid !== changes.node.previousValue.uuid
    ) {
      this.createUppyInstance(changes.node.currentValue);
    }
  }

  ngAfterViewInit() {
    // console.log('[UploadComponent:ngAfterViewInit] ', {
    //   node: this.node,
    //   uuid: this.uuid
    // });
  }

  ngAfterViewChecked() {
    // console.log('[UploadComponent:ngAfterViewInit] ', {
    //   node: this.node,
    //   uuid: this.uuid
    // });
  }

  private createUppyInstance(node: Node) {
    console.log('[UploadComponent:ngAfterViewChecked] ', {
      node: node,
      uuid: this.uuid
    });

    if (!node.uuid) {
      return;
    }

    this.uppyInstance = this.uploadService.configureField(node, this.uuid);

    uppyEvents.forEach(ev =>
      this.uppyInstance.on(ev, (data1, data2, data3) => {
        console.log('[UploadComponent:UppyEvent] ', {
          event: ev,
          data1,
          data2,
          data3
        });

        if (!this.on) {
          return;
        }

        this.on.next([ev, data1, data2, data3]);
      })
    );

    this.uppyInstance.on('complete', () => {
      setTimeout(() => {
        if (this.uppyInstance) {
          this.uppyInstance.reset();
        }
      }, 1500);
    });
  }
}
