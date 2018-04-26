import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../components/common/api';

import { VolumeService } from './volume.service';

@Component({
    selector: 'volume-list',
    templateUrl: 'volumeList.html',
    styleUrls: [],
    animations: []
})
export class VolumeListComponent implements OnInit {
    createSnapshotDisplay = false;
    selectedVolumes = [];
    volumes = [];
    menuItems: MenuItem[];
    volumeForSnapshot;
    label = {};
    snapshotFormGroup;

    constructor(
        // private I18N: I18NService,
        // private router: Router
        private VolumeService: VolumeService,
        private fb: FormBuilder
    ) {
        this.snapshotFormGroup = this.fb.group({
            "name": ["", Validators.required ],
            "description":["", Validators.maxLength(200)]
        })
    }

    ngOnInit() {
        // this.volumes = [
        //     { "name": "vol-01", "status": "Available", "capacity": "200.000 GB", "profile": "PF_block_01", "az":"az_01" },
        //     { "name": "vol-02", "status": "Error", "capacity": "200.000 GB", "profile": "PF_block_02", "az":"az_02" }
        // ];
        this.menuItems = [
            { "label": "Modify", command:()=>{} },
            { "label": "Expand", command:()=>{} },
            { "label": "Delete", command:()=>{} }
        ];

        this.getVolumes();

        this.label = {
            volume: 'Volume',
            name: 'Name',
            description: 'Description'
        }
    }

    getVolumes(){
        this.VolumeService.getVolumes().subscribe((res) => {
            this.volumes = res.json();
        });
    }

    deleteVolume(){
        this.selectedVolumes.forEach(volume => {
            this.VolumeService.deleteVolume(volume.id).subscribe((res) => {
                this.getVolumes();
            });
        });
    }

    showCreateSnapshotDialog(selectVolume){
        this.createSnapshotDisplay = true;
        this.volumeForSnapshot = selectVolume;

    }

    createSnapshot(){
        let param = {
            name: this.snapshotFormGroup.value.name,
            volumeId: this.volumeForSnapshot.id,
            description: this.snapshotFormGroup.value.description
        }
        this.VolumeService.createSnapshot(param).subscribe((res) => {
            this.createSnapshotDisplay = false;
            console.log(res.json());
        });
    }





}
