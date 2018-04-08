import { Component, Input, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../../components/common/api';

@Component({
    selector: 'tenant-detail',
    templateUrl: 'tenantDetail.html',
    styleUrls: ['tenantDetail.scss'],
    animations: []
})
export class TenantDetailComponent implements OnInit {
    users = [];
    @Input() isDetailFinished: Boolean;

    constructor(
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.users = [
            { "username": "admin", "status": "Enabled", "role": "System Administrator" },
            { "username": "cloud_admin", "status": "Enabled", "role": "Storage Administrator" },
            { "username": "jack", "status": "Enabled", "role": "Tenant User" }
        ];

        setTimeout(()=> {
            this.isDetailFinished = true;
        }, 1000);
    }

}
