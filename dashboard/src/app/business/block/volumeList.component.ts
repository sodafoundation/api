import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Router } from '@angular/router';
import { I18NService } from 'app/shared/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem ,ConfirmationService} from '../../components/common/api';

import { VolumeService, SnapshotService,ReplicationService} from './volume.service';
import { ProfileService } from './../profile/profile.service';
import { identifierModuleUrl } from '@angular/compiler';

let _ = require("underscore");
@Component({
    selector: 'volume-list',
    templateUrl: 'volumeList.html',
    styleUrls: [],
    animations: []
})
export class VolumeListComponent implements OnInit {
    createSnapshotDisplay = false;
    createReplicationDisplay = false;
    expandDisplay = false;
    modifyDisplay = false;
    selectVolumeSize :number;
    unit:number = 1;
    repPeriod : number=60;
    capacityOptions = [
        {
            label: 'GB',
            value: 'gb'
        },
        {
            label: 'TB',
            value: 'tb'
        }

    ];
    profileOptions = [
        {
            label: 'Select Profile',
            value: null
        }
    ];
    azOption=[{label:"default",value:"default"}];
    selectedVolumes = [];
    volumes = [];
    menuItems: MenuItem[];
    label = {
        name: this.I18N.keyID['sds_block_volume_name'],
        volume:  this.I18N.keyID['sds_block_volume_title'],
        description:  this.I18N.keyID['sds_block_volume_descri']
    };
    snapshotFormGroup;
    modifyFormGroup;
    expandFormGroup;
    replicationGroup;
    errorMessage = {
        "name": { required: "Name is required." },
        "description": { maxlength: "Max. length is 200." },
        "repName":{ required: "Name is required." },
        "profileOption":{ required: "Name is required." },
        "expandSize":{required: "Expand Capacity is required."}
    };
    profiles;
    selectedVolume;

    constructor(
        public I18N: I18NService,
        private router: Router,
        private VolumeService: VolumeService,
        private SnapshotService: SnapshotService,
        private ProfileService: ProfileService,
        private ReplicationService: ReplicationService,
        private confirmationService: ConfirmationService,
        private fb: FormBuilder
    ) {
        this.snapshotFormGroup = this.fb.group({
            "name": ["", Validators.required],
            "description": ["", Validators.maxLength(200)]
        });
        this.modifyFormGroup = this.fb.group({
            "name": ['', Validators.required]
        });
        this.expandFormGroup = this.fb.group({
            "expandSize":[1,{validators:[Validators.required], updateOn:'change'} ],
            "capacityOption":[this.capacityOptions[0] ]
        });
        this.expandFormGroup.get("expandSize").valueChanges.subscribe(
            (value:string)=>{
                this.selectVolumeSize = parseInt(this.selectedVolume.size) + parseInt(value)*this.unit;
            }
        );
        this.expandFormGroup.get("capacityOption").valueChanges.subscribe(
            (value:string)=>{
                this.unit =(value === "tb" ? 1024: 1);
                this.selectVolumeSize = parseInt(this.selectedVolume.size) + parseInt(this.expandFormGroup.value.expandSize)*this.unit;
            }
        )
        this.replicationGroup = this.fb.group({
            "repName": ['',{validators:[Validators.required], updateOn:'change'}],
            "az": [this.azOption[0]],
            "profileOption":['',{validators:[Validators.required], updateOn:'change'}]
        });

    }

    ngOnInit() {
        this.menuItems = [
            {
                "label": this.I18N.keyID['sds_block_volume_modify'],
                command: () => {
                    this.modifyDisplay = true;
                }
            },
            {
                "label": this.I18N.keyID['sds_block_volume_expand'],
                command: () => {
                    this.expandDisplay = true;
                    this.expandFormGroup.reset();
                    this.selectVolumeSize = parseInt(this.selectedVolume.size);
                    this.expandFormGroup.controls["expandSize"].setValue(1);
                    this.unit = 1;
                }
            },
            {
                "label": this.I18N.keyID['sds_block_volume_delete'], command: () => {
                    if (this.selectedVolume && this.selectedVolume.id) {
                        this.deleteVolumes(this.selectedVolume);
                    }
                }
            }
        ];
        this.getVolumes();

        this.getProfiles()
    }

    getVolumes() {
        this.VolumeService.getVolumes().subscribe((res) => {
            this.volumes = res.json();
            this.volumes.forEach((item)=>
                {
                    this.ProfileService.getProfileById(item.profileId).subscribe((res)=>{
                        item.profileName = res.json().name;
                    })
                }
            )
        });
    }

    getProfiles() {
        this.ProfileService.getProfiles().subscribe((res) => {
            this.profiles = res.json();
            this.profiles.forEach(profile => {
                this.profileOptions.push({
                    label: profile.name,
                    value: profile.id
                });
            });
        });
    }

    batchDeleteVolume() {
        this.selectedVolumes.forEach(volume => {
            this.deleteVolume(volume.id);
        });
    }

    deleteVolumeById(id) {
        this.deleteVolume(id);
    }

    deleteVolume(id) {
        this.VolumeService.deleteVolume(id).subscribe((res) => {
            this.getVolumes();
        });
    }

    createSnapshot() {
        if(!this.snapshotFormGroup.valid){
            for(let i in this.snapshotFormGroup.controls){
                this.snapshotFormGroup.controls[i].markAsTouched();
            }
            return;
        }
        let param = {
            name: this.snapshotFormGroup.value.name,
            volumeId: this.selectedVolume.id,
            description: this.snapshotFormGroup.value.description
        }
        this.SnapshotService.createSnapshot(param).subscribe((res) => {
            this.createSnapshotDisplay = false;
        });
    }

    returnSelectedVolume(selectedVoluem, dialog) {
        if (dialog === 'snapshot') {
            this.createSnapshotDisplay = true;
        } else if (dialog === 'replication') {
            this.createReplicationDisplay = true;
        }
        this.selectedVolume = selectedVoluem;
        this.replicationGroup.reset();
        this.replicationGroup.controls["repName"].setValue(this.selectedVolume.name+"-replication");
        this.replicationGroup.controls["az"].setValue(this.azOption[0]);
        this.selectVolumeSize = parseInt(this.selectedVolume.size) + parseInt(this.expandFormGroup.value.expandSize);
    }

    modifyVolume() {
        let param = {
            name: this.modifyFormGroup.value.name
        };
        this.VolumeService.modifyVolume(this.selectedVolume.id, param).subscribe((res) => {
            this.getVolumes();
            this.modifyDisplay = false;
        });
    }
    expandVolume(){
        if(!this.expandFormGroup.valid){
            for(let i in this.expandFormGroup.controls){
                this.expandFormGroup.controls[i].markAsTouched();
            }
            return;
        }
        let param = {
            "extend": {
                "newSize": this.selectVolumeSize
            }
        }
        this.VolumeService.expandVolume(this.selectedVolume.id, param).subscribe((res) => {
            this.getVolumes();
            this.expandDisplay = false;
        });
    }
    createReplication(){
        if(!this.replicationGroup.valid){
            for(let i in this.replicationGroup.controls){
                this.replicationGroup.controls[i].markAsTouched();
            }
            return;
        }
        let param = {
            "name":this.replicationGroup.value.repName ,
            "size": this.selectedVolume.size,
            "availabilityZone": this.replicationGroup.value.az.value,
            "profileId": this.replicationGroup.value.profileOption,
        }
        this.VolumeService.createVolume(param).subscribe((res) => {
            let param = {
                "name":this.replicationGroup.value.repName ,
                "primaryVolumeId": this.selectedVolume.id,
                "availabilityZone": this.replicationGroup.value.az.value,
                "profileId": this.replicationGroup.value.profileOption,
                "replicationMode":"async",
                "replicationPeriod":Number(this.repPeriod),
                "secondaryVolumeId":res.json().id
            }
            this.ReplicationService.createReplication(param).subscribe((res) => {
                this.getVolumes();
                this.createReplicationDisplay = false;
            });
        });
    }
    deleteVolumes(volumes){
        let arr=[], msg;
        if(_.isArray(volumes)){
            volumes.forEach((item,index)=> {
                arr.push(item.id);
            })
            msg = "<div>Are you sure you want to delete the selected volumes?</div><h3>[ "+ volumes.length +" Volumes ]</h3>";
        }else{
            arr.push(volumes.id);
            msg = "<div>Are you sure you want to delete the volume?</div><h3>[ "+ volumes.name +" ]</h3>";
        }

        this.confirmationService.confirm({
            message: msg,
            header: this.I18N.keyID['sds_block_volume_deleVolu'],
            acceptLabel: this.I18N.keyID['sds_block_volume_delete'],
            isWarning: true,
            accept: ()=>{
                arr.forEach((item,index)=> {
                    this.deleteVolume(item)
                })

            },
            reject:()=>{}
        })

    }
}
