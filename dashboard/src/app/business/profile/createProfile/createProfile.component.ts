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
    showCustomization = false;
    showStoragePool = false;
    msgs: Message[] = [];
    userform: FormGroup;
    submitted: boolean;
    genders: SelectItem[];
    description: string;

    profileform: FormGroup;
    qosPolicy:FormGroup;
    repPolicy:FormGroup;
    snapPolicy:FormGroup;
    paramData= {
        extras:{protocol:""},
        storageType:""
    };
    label;
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
        }/*,
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
        }*/
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
    errorMessage={
        "name": { required: this.I18N.keyID['sds_profile_create_name_require']},
        "maxIOPS": { required: this.I18N.keyID['sds_required'].replace("{0}","MaxIOPS")},
        "maxBWS" :{ required: this.I18N.keyID['sds_required'].replace("{0}","MaxBWS")},
        "repPeriod" :{ required: this.I18N.keyID['sds_required'].replace("{0}","Period")},
        "repBandwidth" :{ required: this.I18N.keyID['sds_required'].replace("{0}","Bandwidth")},
        "repRPO" :{ required: this.I18N.keyID['sds_required'].replace("{0}","RPO")},
        "datetime" :{ required: this.I18N.keyID['sds_required'].replace("{0}","Execution Time")},
        "snapNum" :{ required: this.I18N.keyID['sds_required'].replace("{0}","Retention")},
        "duration" :{ required: this.I18N.keyID['sds_required'].replace("{0}","Retention")},
    };
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
    snapSchedule = [
        {
            label: 'Hourly',
            value: 'Hourly'
        },
        {
            label: 'Daily',
            value: 'Daily'
        },
        {
            label: 'Weekly',
            value: 'Weekly'
        },
        {
            label: 'Monthly',
            value: 'Monthly'
        }
    ];

    weekDays;

    constructor(
        public I18N: I18NService,
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
            MBPS: 'MaxBWS',
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
            'storageType': new FormControl('Thin', Validators.required),
            'policys': new FormControl(''),
            'snapshotRetention': new FormControl('Time')
        });
        this.qosPolicy = this.fb.group({
            "maxIOPS": new FormControl(6000, Validators.required),
            "maxBWS" : new FormControl(100, Validators.required),
        });
        this.repPolicy = this.fb.group({
            "repType": new FormControl("mirror", Validators.required),
            "repMode": new FormControl(this.replicationModeOptions[0].value, Validators.required),
            "repPeriod": new FormControl(60, Validators.required),
            "repBandwidth": new FormControl(10, Validators.required),
            "repRGO": new FormControl(this.replicationRGOOptions[0].value, Validators.required),
            "repRTO": new FormControl(this.replicationRTOOptions[0].value, Validators.required),
            "repRPO": new FormControl(0, Validators.required),
            "repCons": new FormControl([])
        });
        this.snapPolicy = this.fb.group({
            "Schedule": new FormControl(this.snapSchedule[0].value, Validators.required),
            "datetime": new FormControl("00:00", Validators.required),
            "snapNum": new FormControl(1, Validators.required),
            "duration": new FormControl(0, Validators.required),
            "retentionOptions": new FormControl(this.snapshotRetentionOptions[0].value)
        });
        this.paramData= {
            extras:{protocol:this.profileform.value.protocol},
            storageType:this.profileform.value.storageType
        };
        this.profileform.get("protocol").valueChanges.subscribe(
            (value:string)=>{
                this.paramData = {
                    extras:{protocol:value},
                    storageType:this.profileform.value.storageType
                }
            }
        );
        this.profileform.get("storageType").valueChanges.subscribe(
            (value:string)=>{
                this.paramData = {
                    extras:{protocol:this.profileform.value.protocol},
                    storageType:value
                }
            }
        );



    }

    onSubmit(value) {
        this.submitted = true;
        this.msgs = [];
        this.msgs.push({ severity: 'info', summary: 'Success', detail: 'Form Submitted' });
        this.param.name = value.name;
        this.param.storageType = value.storageType;
        this.param.extras.protocol = value.protocol;
        this.param.extras.policys = value.policys;
        if(this.qosIsChecked){
            if(!this.qosPolicy.valid){
                for(let i in this.qosPolicy.controls){
                    this.qosPolicy.controls[i].markAsTouched();
                }
                return;
            }
            this.param.extras[":provisionPolicy"]= {
                "ioConnectivityLoS": {
                    "maxIOPS": this.qosPolicy.value.maxIOPS,
                    "maxBWS": this.qosPolicy.value.maxBWS
                }
            }
        }
        if(this.replicationIsChecked){
            if(!this.repPolicy.valid){
                for(let i in this.repPolicy.controls){
                    this.repPolicy.controls[i].markAsTouched();
                }
                return;
            }
            this.param.extras["replicationType"]= "ArrayBased";
            this.param.extras[":replicationPolicy"]={
                "dataProtectionLoS": {
                    "replicaTypes": this.repPolicy.value.repType,
                    "recoveryGeographicObject": this.repPolicy.value.repRGO,
                    "recoveryPointObjective": this.repPolicy.value.repRPO,
                    "recoveryTimeObjective": this.repPolicy.value.repRTO,
                },
                "replicaInfos": {
                    "replicaUpdateMode": this.repPolicy.value.repMode,
                    "consistencyEnabled": this.repPolicy.value.repCons.length==0 ? false:true,
                    "replicationPeriod": this.repPolicy.value.repPeriod,
                    "replicationBandwidth": this.repPolicy.value.repBandwidth
                }
            }
        }
        if(this.snapPolicy){
            if(!this.snapPolicy.valid){
                for(let i in this.snapPolicy.controls){
                    this.snapPolicy.controls[i].markAsTouched();
                }
                return;
            }
            let reten = this.snapPolicy.value.retentionOptions === "Quantity" ? {
                    "number": this.snapPolicy.value.snapNum,
                }:{"duration": this.snapPolicy.value.duration}
            this.param.extras[":snapshotPolicy"]= {
                "schedule": {
                    "datetime": "1970-01-01T"+this.snapPolicy.value.datetime+":00",
                    "occurrence": this.snapPolicy.value.Schedule //Monthly, Weekly, Daily, Hourly
                },
                "retention": reten
            }
        }
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
