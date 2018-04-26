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

    unspecified: boolean = false;

    ownedTenant: string = "";

    constructor(
        private http: Http,
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.userID = this.detailUserInfo;

        let request: any = { params: {} };
        this.http.get("/v3/users/" + this.userID + "/projects", request).subscribe((res) => {
            if(res.json().projects.length == 0){
                this.ownedTenant = "Unspecified";
                this.unspecified = true;
            }
            res.json().projects.forEach((ele, i) => {
                if(i==0){
                    this.ownedTenant = ele.name;
                }else{
                    this.ownedTenant += ", "+ ele.name;
                }
            });

            this.isUserDetailFinished = true;
        });

    }

}
