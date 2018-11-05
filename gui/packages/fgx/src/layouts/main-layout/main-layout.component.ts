import { Component, OnInit } from '@angular/core';

import { NbMenuItem, NbMenuService, NbSidebarService } from '@nebular/theme';

@Component({
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss']
})
export class MainLayoutComponent implements OnInit {
  menuItems: NbMenuItem[];
  position: string;

  constructor(
    private sidebarService: NbSidebarService,
    private menuService: NbMenuService
  ) {}

  ngOnInit() {
    this.position = 'normal';
    this.menuItems = [
      {
        title: 'Home',
        link: '/dashboard',
        icon: 'nb-home',
        home: true
      },
      {
        title: 'Content',
        icon: 'nb-keypad',
        children: [
          {
            title: 'Media',
            link: '/media'
          }
        ]
      },
      {
        title: 'Sites',
        icon: 'nb-tables'
      },
      {
        title: 'Configuration',
        icon: 'nb-gear',
        children: [
          {
            title: 'Logout',
            link: '/guard/logout'
          }
        ]
      }
    ];
  }

  goToHome() {
    this.menuService.navigateHome();
  }

  toggleSidebar(): boolean {
    this.sidebarService.toggle(true, 'menu-sidebar');
    return false;
  }
}
