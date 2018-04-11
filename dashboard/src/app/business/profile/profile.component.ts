import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate} from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

@Component({
    templateUrl: './profile.component.html',
    styleUrls: [
        './profile.component.css'
    ],
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
export class ProfileComponent implements OnInit{
    profiles = [];
    constructor(
        // private I18N: I18NService,
        // private router: Router
    ){}
    showCard = true;
    ngOnInit() {
        this.profiles = [
            {
              "id": "5d8c3732-a248-50ed-bebc-539a6ffd25c1",
              "name": "Gold",
              "protocol": "FC",
              "type": "Thin",
              "policys": [
                  "Qos",
                  "Snapshot",
                  "Replication"
              ],
              "description": "provide gold storage service",
              "extras": {
                "key1": "value1",
                "key2": {
                  "subKey1": "subValue1",
                  "subKey2": "subValue2"
                },
                "key3": "value3"
              }     
            },
            {
              "id": "5d8c3732-a248-50ed-bebc-539a6ffd25c2",
              "name": "Silver",
              "protocol": "iSCSI",
              "type": "Thick",
              "policys": [
                "Qos",
                "Snapshot"
            ],
              "description": "provide silver storage service",
              "extras": {
                "key1": "value1",
                "key2": "value2"
              }
            }
          ];
    }
    onDeleteProfile = function(){
        alert(1);
    }
}
