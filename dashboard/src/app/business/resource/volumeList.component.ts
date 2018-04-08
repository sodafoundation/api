import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { MenuItem } from '../../components/common/api';

@Component({
    selector: 'volume-list',
    templateUrl: 'volumeList.html',
    styleUrls: [],
    animations: []
})
export class VolumeListComponent implements OnInit {
    volumes = [];
    menuItems: MenuItem[];

    constructor(
        // private I18N: I18NService,
        // private router: Router
    ) { }

    ngOnInit() {
        this.volumes = [
            { "name": "vol-01", "status": "Available", "capacity": "200.000 GB", "profile": "PF_block_01", "az":"az_01" },
            { "name": "vol-02", "status": "Error", "capacity": "200.000 GB", "profile": "PF_block_02", "az":"az_02" }
        ];
        this.menuItems = [
            { "label": "Modify", command:()=>{} },
            { "label": "Expand", command:()=>{} },
            { "label": "Delete", command:()=>{} }
        ];
    }

}
