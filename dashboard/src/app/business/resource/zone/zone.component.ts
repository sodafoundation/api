import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService, ParamStorService} from 'app/shared/api';
import { Http } from '@angular/http';
import { trigger, state, style, transition, animate } from '@angular/animations';

@Component({
    selector: 'zone-table',
    templateUrl: './zone.html',
    styleUrls: [],
    animations: [
        trigger('overlayState', [
            state('hidden', style({
                opacity: 0
            })),
            state('visible', style({
                opacity: 1
            })),
            transition('visible => hidden', animate('400ms ease-in')),
            transition('hidden => visible', animate('400ms ease-out'))
        ]),

        trigger('notificationTopbar', [
            state('hidden', style({
                height: '0',
                opacity: 0
            })),
            state('visible', style({
                height: '*',
                opacity: 1
            })),
            transition('visible => hidden', animate('400ms ease-in')),
            transition('hidden => visible', animate('400ms ease-out'))
        ])
    ]
})
export class ZoneComponent implements OnInit {

    zones = [];

    constructor(
        public I18N: I18NService,
        private http: Http,
        private paramStor: ParamStorService
    ) { }

    ngOnInit() {

        this.listAZ();
    }

    listAZ(){
        this.zones = [];
        let reqUser: any = { params:{} };
        let user_id = this.paramStor.CURRENT_USER().split("|")[1];
        this.http.get("/v3/users/"+ user_id +"/projects", reqUser).subscribe((objRES) => {
            let project_id;
            objRES.json().projects.forEach(element => {
                if(element.name == "admin"){
                    project_id = element.id;
                }
            })

            let reqPool: any = { params:{} };
            this.http.get("/v1beta/"+ project_id +"/pools", reqPool).subscribe((poolRES) => {
                let AZs=[];
                poolRES.json().forEach(ele => {
                    if(!AZs.includes(ele.availabilityZone)){
                        AZs.push(ele.availabilityZone);

                        let [name,region,description] = [ele.availabilityZone, "default_region", "--"];
                        this.zones.push({name,region,description});
                    }
                })
                console.log(this.zones);
            })

        })
    }
}

