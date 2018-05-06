import { Component, OnInit, Input } from '@angular/core';
import { FormControl, FormGroup, AbstractControl } from '@angular/forms';

@Component({
  selector: 'gonode-input-error',
  templateUrl: './input-error.component.html',
  styleUrls: ['./input-error.component.css']
})
export class InputErrorComponent implements OnInit {
  @Input() name: string;
  @Input() form: FormGroup;
  field: AbstractControl;

  constructor() {}

  ngOnInit() {
    this.field = this.form.get(this.name);
  }
}
