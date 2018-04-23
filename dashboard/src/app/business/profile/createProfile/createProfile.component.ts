import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Validators,FormControl,FormGroup,FormBuilder } from '@angular/forms';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { I18nPluralPipe } from '@angular/common';

import { Message,SelectItem } from './../../../components/common/api';
import { ProfileService } from './../profile.service';

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
    showCustomization = false;
    msgs:Message[] = [];
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
            label:'Please Select Protocol',
            value:''
        },
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

    //用户自定义key、value，用于双向数据绑定
    customizationKey = '';
    customizationValue = '';

    //用户自定义项，用于
    customizationItems = [];


    constructor(
        // private I18N: I18NService,
        private router: Router,
        private ProfileService: ProfileService,
        private fb: FormBuilder
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
            value: 'Value',
            maxIOPS: 'MaxIOPS',
            MBPS: 'MBPS'
        };

        this.profileform = this.fb.group({
            'name': new FormControl('', Validators.required),
            'protocol': new FormControl('', Validators.required),
            'storageType': new FormControl('', Validators.required),
            'policys': new FormControl('')
        });
    }

    onSubmit(value) {
        this.submitted = true;
        this.msgs = [];
        this.msgs.push({severity:'info', summary:'Success', detail:'Form Submitted'});
        this.param.name = value.name;
        this.param.storageType = value.storageType;
        this.param.extras.protocol = value.protocol;
        this.param.extras.policys = value.policys;
        if(this.customizationItems.length > 0){
            let arrLength = this.customizationItems.length;
            for(let i=0;i<arrLength;i++){
                this.param.extras[this.customizationItems[i].key] = this.customizationItems[i].value;
            }
        }
        this.createProfile(this.param);
    }

    createProfile(param) {
        this.ProfileService.createProfile(param).subscribe((res) => {
            // return res.json();
            // this.profiles = res.json();
            this.router.navigate(['/profile']);
        })
    }

    getI18n() {

        // return {};
    }

    showDetails(policyType){
        // alert(policyType);
        this[policyType+'IsChecked'] = !this[policyType+'IsChecked'];
    }

    addCustomization(){
        this.customizationItems.push({
            key:this.customizationKey,
            value:this.customizationValue
        });
        // this.param.extras[this.customizationKey] = this.customizationValue;
        this.showCustomization = false
        this.customizationKey = '';
        this.customizationValue = '';
        console.log(this.customizationItems);
    }

    deleteCustomization(index){
        this.customizationItems.splice(index, 1);
        console.log(this.customizationItems);
    }

}
