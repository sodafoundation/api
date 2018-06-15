import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from './../../../../app/shared/api';
import { AppService } from './../../../../app/app.service';
import { trigger, state, style, transition, animate} from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

@Component({
    selector: 'region-table',
    templateUrl: './region.html',
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
export class RegionComponent implements OnInit{
    regions = [];

    constructor(
         public I18N: I18NService,
        // private router: Router
    ){}

    ngOnInit() {
        this.regions = [
            { "name": this.I18N.keyID['sds_resource_region_default'], "role": "Primary Region", }
        ];
    }

    rowSelect(rowdata){
        console.log(rowdata);
    }

}

