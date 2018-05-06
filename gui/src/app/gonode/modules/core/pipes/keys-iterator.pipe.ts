import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'keysIterator'
})
export class KeysIteratorPipe implements PipeTransform {
  transform(obj: object) {
    return Object.keys(obj);
  }
}
