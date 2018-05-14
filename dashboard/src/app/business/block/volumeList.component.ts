import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Router } from '@angular/router';
import { I18NService } from 'app/shared/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../components/common/api';

import { VolumeService, SnapshotService } from './volume.service';
import { ProfileService } from './../profile/profile.service';
import { identifierModuleUrl } from '@angular/compiler';

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
    selectedVolumes = [];
    volumes = [];
    menuItems: MenuItem[];
    label = {};
    snapshotFormGroup;
    modifyFormGroup;
    expandFormGroup;
    errorMessage = {
        "name": { required: "Name is required." },
        "description": { maxlength: "Max. length is 200." }
    };
    profiles;
    selectedVolume;

    constructor(
        // private I18N: I18NService,
        private router: Router,
        private VolumeService: VolumeService,
        private SnapshotService: SnapshotService,
        private ProfileService: ProfileService,
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
            "capacityOption":[this.capacityOptions[0],{validators:[Validators.required], updateOn:'change'} ]
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
    }

    ngOnInit() {
        this.menuItems = [
            {
                "label": "Modify",
                command: () => {
                    this.modifyDisplay = true;
                }
            },
            {
                "label": "Expand",
                command: () => {
                    this.expandDisplay = true;
                    this.expandFormGroup.reset();
                    this.selectVolumeSize = parseInt(this.selectedVolume.size);
                    this.expandFormGroup.controls["expandSize"].setValue(1);
                    this.unit = 1;
                }
            },
            {
                "label": "Delete", command: () => {
                    if (this.selectedVolume && this.selectedVolume.id) {
                        this.deleteVolumeById(this.selectedVolume.id);
                    }
                }
            }
        ];
        this.getVolumes();

        this.label = {
            volume: 'Volume',
            name: 'Name',
            description: 'Description'
        }
        this.getProfiles()
    }

    getVolumes() {
        this.VolumeService.getVolumes().subscribe((res) => {
            this.volumes = res.json();
        });
    }

    getProfiles() {
        this.ProfileService.getProfiles().subscribe((res) => {
            this.profiles = res.json();
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
}
