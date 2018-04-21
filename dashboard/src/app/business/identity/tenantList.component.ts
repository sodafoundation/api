import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../components/common/api';

@Component({
    selector: 'tenant-list',
    templateUrl: 'tenantList.html',
    styleUrls: [
        'dialogcss.css'
    ],
    animations: []
})
export class TenantListComponent implements OnInit {
    tenants = [];
    isDetailFinished = false;
    createTenantDisplay = false;
    
    sortField: string;

    constructor(
        private http: Http,
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.listTenants();

    }

    listTenants() {
        this.sortField = "name";

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/projects", request).subscribe((res) => {
            this.tenants = res.json().projects;
        });
    }

    showCreateTenant() {
        this.createTenantDisplay = true;
    }

    onRowExpand(evt) {
        this.isDetailFinished = false;
        console.log(evt.data);

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        // this.http.get("/v2.0/tenants/"+ evt.data.id +"/users").subscribe((res) => {
        //     // this.tenants = res.json().projects;
        // });

    }

    label: object = {
        userNameLabel: 'Name:',
        descriptionLabel: 'Description:',
    }

}






