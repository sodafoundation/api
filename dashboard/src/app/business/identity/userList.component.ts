import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../components/common/api';

@Component({
    selector: 'user-list',
    templateUrl: 'userList.html',
    styleUrls: [],
    animations: []
})
export class UserListComponent implements OnInit {
    users = [];

    constructor(
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.users = [
            { "username": "admin", "status": "Enabled", "tenant": "tenant_A, tenant_B", "role": "System Administrator", "userid":"uu220816001" },
            { "username": "cloud_admin", "status": "Enabled", "tenant": "tenant_B", "role": "Storage Administrator", "userid":"uu220816002" },
            { "username": "jack", "status": "Enabled", "tenant": "tenant_A", "role": "Tenant User", "userid":"uu220816003" }
        ];
    }

}
