import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../components/common/api';

import { FormGroup, FormControl } from '@angular/forms';

@Component({
    selector: 'user-list',
    templateUrl: 'userList.html',
    styleUrls: [
        'dialogcss.css'
    ],
    animations: []
})
export class UserListComponent implements OnInit {
    users = [];
    createUserDisplay = false;
    constructor(
        // private I18N: I18NService,
        // private router: Router
    ) { }
    
    showCreateUser(): void{
        this.createUserDisplay = true;
    }
    

    roles = [
        {label:'Select User', value:null},
        {label:'New York', value:{id:1, name: 'New York', code: 'NY'}},
        {label:'Rome', value:{id:2, name: 'Rome', code: 'RM'}},
        {label:'London', value:{id:3, name: 'London', code: 'LDN'}},
        {label:'Istanbul', value:{id:4, name: 'Istanbul', code: 'IST'}},
        {label:'Paris', value:{id:5, name: 'Paris', code: 'PRS'}}
    ];

    label:object = {
        userNameLabel:'Username:',
        passwordLabel:'Password:',
        confirmPasswordLabel:'Confirm Password:',
        roleLabel:'Role:',
        tenantLabel:'Tenant:'
    }
    tenants = [
        {label:'Select Tenant', value:null},
        {label:'Tenant1', value:1},
        {label:'Tenant2', value:2}
    ]

    errorMessage = {};


    ngOnInit() {
        this.users = [
            { "username": "admin", "status": "Enabled", "tenant": "tenant_A, tenant_B", "role": "System Administrator", "userid":"uu220816001" },
            { "username": "cloud_admin", "status": "Enabled", "tenant": "tenant_B", "role": "Storage Administrator", "userid":"uu220816002" },
            { "username": "jack", "status": "Enabled", "tenant": "tenant_A", "role": "Tenant User", "userid":"uu220816003" }
        ];
    }

}
