import { Component, Input, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../../components/common/api';

@Component({
    selector: 'user-detail',
    templateUrl: 'userDetail.html',
    styleUrls: ['userDetail.scss'],
    animations: []
})
export class userDetailComponent implements OnInit {
    @Input() isUserDetailFinished: Boolean;

    @Input() detailUserInfo: string;

    userID: string;

    userRole: string;

    defaultTenant: string;

    constructor(
        private http: Http,
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.userID = this.detailUserInfo.split("|")[0];
        this.userRole = this.detailUserInfo.split("|")[1] == "admin" ? "System Administrator" : "Tenant User";

        let request: any = { params: {} };
        this.http.get("/v3/users/" + this.userID, request).subscribe((res) => {
            let project_id = res.json().user.default_project_id;
            let req: any = { params: {} };
            this.http.get("/v3/projects/" + project_id, req).subscribe((res) => {
                this.defaultTenant = res.json().project.name;

                this.isUserDetailFinished = true;
            })

        });

    }

}
