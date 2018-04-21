import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener, ViewChildren } from '@angular/core';
import { Http } from '@angular/http';
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
    tenantUsers = [];
    tenantLists = [];
    createUserDisplay = false;
    isUserDetailFinished = false;
    selectedTenant: string;
    username: string;
    description: string = "Description";
    password: string;
    userRole: string;

    detailUserInfo: string;

    sortField: string;
    constructor(
        private http: Http,
        // private I18N: I18NService,
        // private router: Router
    ) { }

    
    label:object = {
        userNameLabel:'Username:',
        passwordLabel:'Password:',
        confirmPasswordLabel:'Confirm Password:',
        roleLabel:'Role:',
        tenantLabel:'Tenant:'
    }

    errorMessage = {};
    
    showCreateUser(): void{
        this.createUserDisplay = true;
        this.getTenants();
        this.getRoles();
    }

    createUser(){
        let request: any = { user:{} };
        request.user = {
            "default_project_id": this.selectedTenant,
            "domain_id": "default",
            "name": this.username,
            "description": this.description,
            "password": this.password
        }
        
        this.http.post("/v3/users", request).subscribe((res) => {
            let userInfo = res.json().user;
            let request: any = { params:{} };
            this.http.put("/v3/projects/"+ userInfo.default_project_id +"/users/"+ userInfo.id +"/roles/"+ this.userRole, request).subscribe((r) => {
                this.createUserDisplay = false;
                this.tenantUsers = [];
                this.listUsers();
            })
        });
    }
    
    getRoles(){
        let request: any = { params:{} };
        this.http.get("/v3/roles", request).subscribe((res) => {
            res.json().roles.forEach((item, index) => {
                if(item.name == "Member"){
                    this.userRole = item.id;
                }
            })
        });
    }
    getTenants(){
        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }

        this.http.get("/v3/projects", request).subscribe((res) => {
            res.json().projects.map((item, index) => {
                let tenant = {};
                tenant["label"] = item.name;
                tenant["value"] = item.id;
                this.tenantLists.push(tenant);
                this.selectedTenant =  this.tenantLists[0].value;
            });
        });
    }


    ngOnInit() {
        this.listUsers();
        
    }

    listUsers(){
        this.sortField = "username";

        let request: any = { params:{} };
        request.params = {
            "domain_id": "default"
        }
        this.http.get("/v3/users", request).subscribe((res) => {
            res.json().users.map((item, index) => {
                let user = {};
                user["status"] = (item.enabled == true) ? "Enabled" : "Disabled";
                user["username"] = item.name;
                user["userid"] = item.id;
                user["defaultTenant"] = item.default_project_id;
                user["description"] = item.description;
                this.tenantUsers.push(user);
            });
            console.log(this.tenantUsers);
        });
    }

    onRowExpand(evt) {
        this.isUserDetailFinished = false;
        this.detailUserInfo = evt.data.userid+ "|"+ evt.data.username;
    }
}
