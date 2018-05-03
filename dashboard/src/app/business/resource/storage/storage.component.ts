import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from './../../../../app/shared/api';
import { AppService } from './../../../../app/app.service';
import { trigger, state, style, transition, animate} from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

@Component({
    selector: 'storage-table',
    templateUrl: './storage.html',
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
export class StorageComponent implements OnInit{

    storages = [];

    constructor(
        // private I18N: I18NService,
        // private router: Router
    ){}
    
    ngOnInit() {
        this.storages = [
            { 
                "name": "Dorado6000V3.SiteA.78",
                "ip": "8.46.191.46",
                "status": "Enabled",
                "vender": "Huawei",
                "model": "Dorado6000 V3",
                "region": "Region_Chengdu",
                "zone":"default"
            },
            { 
                "name": "Dorado6000V3.SiteB.115",
                "ip": "8.46.185.114",
                "status": "Enabled",
                "vender": "Huawei",
                "model": "Dorado6000 V3",
                "region": "Region_Chengdu",
                "zone":"remote"
            }
        ];
    }
    
}

