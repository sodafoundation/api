import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Validators, FormControl, FormGroup, FormBuilder } from '@angular/forms';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

import { Message, SelectItem } from './../../../components/common/api';
import { ProfileService } from './../profile.service';

@Component({
    templateUrl: './createProfile.component.html',
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
export class CreateProfileComponent implements OnInit {
    errorMessage;
    showCustomization = false;
    msgs: Message[] = [];
    userform: FormGroup;
    submitted: boolean;
    genders: SelectItem[];
    description: string;

    profileform: FormGroup;

    label = {};
    param = {
        name: '',
        storageType: '',
        description: '',
        extras: {
            protocol: 'iSCSI',
            policys: []
        }
    };
    qosIsChecked = false;
    replicationIsChecked = false;
    snapshotIsChecked = false;
    protocolOptions = [
        {
            label: 'iSCSI',
            value: 'iSCSI'
        },
        {
            label: 'FC',
            value: 'FC'
        },
        {
            label: 'RBD',
            value: 'RBD'
        }
    ];

    //用户自定义key、value，用于双向数据绑定
    customizationKey = '';
    customizationValue = '';

    //用户自定义项，用于
    customizationItems = [];

    replicationTypeOptions = [
        {
            label: 'Mirror',
            value: 'mirror'
        },
        {
            label: 'Snapshot',
            value: 'snapshot'
        },
        {
            label: 'Clone',
            value: 'clone'
        },
        {
            label: 'Tokenized Clone',
            value: 'tokenized'
        }
    ];

    replicationRGOOptions = [
        {
            label: 'Availability Zone',
            value: 'availabilityZone'
        },
        {
            label: 'Rack',
            value: 'rack'
        },
        {
            label: 'Row',
            value: 'row'
        },
        {
            label: 'Server',
            value: 'server'
        },
        {
            label: 'Facility',
            value: 'facility'
        },
        {
            label: 'Region',
            value: 'region'
        }
    ];

    replicationModeOptions = [
        {
            label: 'Synchronous',
            value: 'Synchronous'
        },
        {
            label: 'Asynchronous',
            value: 'Asynchronous'
        },
        {
            label: 'Active',
            value: 'Active'
        },
        {
            label: 'Adaptive',
            value: 'Adaptive'
        }
    ];

    replicationRTOOptions = [
        {
            label: 'Immediate',
            value: 'Immediate'
        },
        {
            label: 'Online',
            value: 'Online'
        },
        {
            label: 'Nearline',
            value: 'Nearline'
        },
        {
            label: 'Offline',
            value: 'Offline'
        }
    ];

    replicationRPOOptions = [
        {
            label: '0',
            value: 0
        },
        {
            label: '4',
            value: 4
        },
        {
            label: '60',
            value: 60
        },
        {
            label: '3600',
            value: 3600
        }
    ];

    snapshotRetentionOptions = [
        {
            label: 'Time',
            value: 'Time'
        },
        {
            label: 'Quantity',
            value: 'Quantity'
        }
    ];

    weekDays;

    constructor(
        // private I18N: I18NService,
        private router: Router,
        private ProfileService: ProfileService,
        private fb: FormBuilder
    ) {
        this.weekDays = [
            {
                label: 'Sun',
                value: 0,
                styleClass: 'week-select-list'
            },
            {
                label: 'Mon',
                value: 1,
                styleClass: 'week-select-list'
            },
            {
                label: 'Tue',
                value: 2,
                styleClass: 'week-select-list'
            },
            {
                label: 'Wed',
                value: 3,
                styleClass: 'week-select-list'
            },
            {
                label: 'Thu',
                value: 4,
                styleClass: 'week-select-list'
            },
            {
                label: 'Fri', value: 5,
                styleClass: 'week-select-list'
            },
            {
                label: 'Sat',
                value: 6,
                styleClass: 'week-select-list'
            }
        ];
    }

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
            value: 'Value',
            maxIOPS: 'MaxIOPS',
            MBPS: 'MBWS',
            replicationLabel: {
                type: 'Type',
                RGO: 'RGO',
                Mode: 'Mode',
                RTO: 'RTO',
                Period: 'Period',
                RPO: 'RPO',
                Bandwidth: 'Bandwidth',
                Consistency: 'Consistency'
            },
            snapshotLabel: {
                Schedule: 'Schedule',
                executionTime: 'Execution Time',
                Retention: 'Retention'
            }
        };

        this.profileform = this.fb.group({
            'name': new FormControl('', Validators.required),
            'protocol': new FormControl('iSCSI'),
            'storageType': new FormControl('', Validators.required),
            'policys': new FormControl(''),
            'snapshotRetention': new FormControl('Time')
        });



    }

    onSubmit(value) {
        this.submitted = true;
        this.msgs = [];
        this.msgs.push({ severity: 'info', summary: 'Success', detail: 'Form Submitted' });
        this.param.name = value.name;
        this.param.storageType = value.storageType;
        this.param.extras.protocol = value.protocol;
        this.param.extras.policys = value.policys;
        if (this.customizationItems.length > 0) {
            let arrLength = this.customizationItems.length;
            for (let i = 0; i < arrLength; i++) {
                this.param.extras[this.customizationItems[i].key] = this.customizationItems[i].value;
            }
        }
        this.createProfile(this.param);
    }

    retentionChange(){
        console.log(this.profileform.controls['snapshotRetention'].value);
    }

    createProfile(param) {
        this.ProfileService.createProfile(param).subscribe((res) => {
            // return res.json();
            // this.profiles = res.json();
            this.router.navigate(['/profile']);
        });
    }



    getI18n() {
        // return {};
    }

    showDetails(policyType) {
        this[policyType + 'IsChecked'] = !this[policyType + 'IsChecked'];
    }

    addCustomization() {
        this.customizationItems.push({
            key: this.customizationKey,
            value: this.customizationValue
        });
        this.showCustomization = false
        this.customizationKey = '';
        this.customizationValue = '';
    }

    deleteCustomization(index) {
        this.customizationItems.splice(index, 1);
    }

}
