import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService, ParamStorService} from 'app/shared/api';
import { Http } from '@angular/http';
import { trigger, state, style, transition, animate } from '@angular/animations';
import {AvailabilityZonesService} from '../resource.service';

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
        private paramStor: ParamStorService,
        private availabilityZonesService: AvailabilityZonesService
    ) { }

    ngOnInit() {

        this.listAZ();
    }

    listAZ(){
        this.zones = [];
        this.availabilityZonesService.getAZ().subscribe((azRes) => {
            let AZs=azRes.json();
            if(AZs && AZs.length !== 0){
                AZs.forEach(item =>{
                    let [name,region,description] = [item, "default_region", "--"];
                    this.zones.push({name,region,description});
                })
            }
        })
    }
}

