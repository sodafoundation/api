import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

import { ProfileService } from './profile.service'

@Component({
    templateUrl: './profile.component.html',
    styleUrls: [
        
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
export class ProfileComponent implements OnInit {
    profileId;
    profiles;
    showWarningDialog = false;
    constructor(
        // private I18N: I18NService,
        // private router: Router
        private ProfileService: ProfileService
    ) { }
    showCard = true;
    ngOnInit() {
        this.getProfiles();

        this.profiles = [
            // {
            //     "id": "bbb",
            //     "name": "Gold",
            //     "protocol": "FC",
            //     "type": "Thin",
            //     "policys": [
            //         "Qos",
            //         "Snapshot",
            //         "Replication"
            //     ],
            //     "description": "provide gold storage service",
            //     "extras": {
            //         "key1": "value1",
            //         "key2": {
            //             "subKey1": "subValue1",
            //             "subKey2": "subValue2"
            //         },
            //         "key3": "value3"
            //     }
            // }
        ];
    }

    getProfiles() {
        this.ProfileService.getProfiles().subscribe((res) => {
            // return res.json();
            this.profiles = res.json();
        });
    }

    showWarningDialogFun(id) {
        this.profileId = id;
        this.showWarningDialog = true;
    }
    deleteProfile(id) {
        this.ProfileService.deleteProfile(id).subscribe((res) => {
            // this.profiles = res.json();
            this.getProfiles();
            this.showWarningDialog = false;
        });
    }
}
