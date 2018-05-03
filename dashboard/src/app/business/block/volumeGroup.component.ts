import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate} from '@angular/animations';

@Component({
    selector: 'volume-group-list',
    templateUrl: 'volumeGroup.html',
    styleUrls: [],
    animations: []
})
export class VolumeGroupComponent implements OnInit{
    volumeGroups=[];
    selectedVolumeGroups=[];

    constructor(
        // private I18N: I18NService,
        // private router: Router
    ){}
    
    ngOnInit() {
      this.volumeGroups = [
          {"name": "group_for_REP", "status": "Available", "profile": "PF_block_01", "volumes": "2"},
          {"name": "group_app01", "status": "Error", "profile": "PF_block_02", "volumes": "5"}
      ]
    }
    
}
