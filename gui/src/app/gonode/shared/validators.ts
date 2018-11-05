import { ValidatorFn, AbstractControl } from '@angular/forms';
import { NodeStatus } from './services/api.service';

export function nodeStatusValidator(): ValidatorFn {
  return (control: AbstractControl): { [key: string]: any } => {
    const status = control.value as number;

    let found = false;

    NodeStatus.forEach(v => {
      if (v[0] === status) {
        found = true;
      }
    });

    if (!found) {
      return { nodeStatus: { value: status } };
    }

    return null;
  };
}
