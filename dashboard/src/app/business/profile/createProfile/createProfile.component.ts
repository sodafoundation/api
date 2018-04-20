import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

@Component({
    templateUrl: './createProfile.component.html',
    styleUrls: [
        './createProfile.component.scss'
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
export class CreateProfileComponent implements OnInit {
    errorMessage;
    label = {};
    param = {
        name: '',
        storageType: '',
        description: '',
        extras: {
            protocol: '',
            policy: []
        }
    };
    qosIsChecked = false;
    replicationIsChecked = false;
    snapshotIsChecked = false;
    protocolOptions = [
        {
            label:'iSCSI',
            value:'iSCSI'
        },
        {
            label:'FC',
            value:'FC'
        },
        {
            label:'RBD',
            value:'RBD'
        }
    ];
    constructor(
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.label = {
            name: 'Name',
            protocol: 'Access Protocol',
            type: 'Provisioning Type',
            qosPolicy: 'QoS Policy',
            replicationPolicy: 'Replication Policy',
            snapshotPolicy: 'Snapshot Policy',
            customization: 'Customization',
            storagePool: 'Available Storage Pool',
            key: 'Key',
            value: 'Value'
        };
    }

    createProfile() {
        alert(1);
    }

    getI18n() {

        // return {};
    }

    showDetails(policyType){
        // alert(policyType);
        this[policyType+'IsChecked'] = !this[policyType+'IsChecked'];
    }
}
