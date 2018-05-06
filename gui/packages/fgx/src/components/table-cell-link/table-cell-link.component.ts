import { Component, OnInit, Input } from '@angular/core';

export class Link {
  href: string;
  title: string;
  target: string;
  text: string;
}

@Component({
  selector: 'fgx-cell-router-link',
  template: `
      <a *ngIf="link.target === 'internal'" routerLink="{{ link.href }}" title="{{ link.title}}">{{ link.text }}</a>
      <a *ngIf="link.target !== 'internal'" href="{{ link.href }}" target="{{ target }}" title="{{ link.title}}">{{ link.text }}</a>
    `
})
export class TableLinkComponent implements OnInit {
  @Input()
  value: string | object;

  @Input()
  rowData: any;

  link: Link;

  ngOnInit() {
    if (typeof this.value === 'string') {
      this.link = {
        href: this.value,
        title: this.value,
        target: 'internal',
        text: this.value
      };
    }

    if (this.value instanceof Link) {
      this.link = this.value;
    }

    if (typeof this.value === 'object') {
      this.link = {
        href: 'href' in this.value ? this.value['href'] : '',
        title: 'title' in this.value ? this.value['title'] : '',
        target: 'target' in this.value ? this.value['target'] : 'internal',
        text: 'text' in this.value ? this.value['text'] : ''
      };
    }
  }
}
