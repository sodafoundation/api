import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate} from '@angular/animations';
import { DialogModule } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';
import { VolumeService ,VolumeGroupService} from './volume.service';

@Component({
    selector: 'volume-group-list',
    templateUrl: 'volumeGroup.html',
    styleUrls: [],
    animations: []
})
export class VolumeGroupComponent implements OnInit{
    volumeGroups=[];
    volemeOptions = [];
    selectedOption :string;
    selectedVolumeGroups=[];
    showVolumeGroupDialog :boolean = false;
    volumeGroupForm:any;
    validRule= {
        'name':'^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){0,127}$'
    };
    constructor(
        // private I18N: I18NService,
        // private router: Router
        private volumeGroupService : VolumeGroupService,
        private fb : FormBuilder
    ){
        this.volumeGroupForm = this.fb.group({
            "group_name":["",{validators:[Validators.required, Validators.pattern(this.validRule.name)], updateOn:'change'} ],
            "profile":[""]
        });
    }
    errorMessage = {
        "group_name": { required: "group name is required.", pattern:"Beginning with a letter with a length of 1-128, it can contain letters / numbers / underlines."},
        "profile": { required: "profile is required."}
    };
    
    label = {
        group_name_lable:'Group Name',
        profile_label:'Profile'
    }

    ngOnInit() {
      this.volumeGroups = [
          {"name": "group_for_REP", "status": "Available", "profile": "PF_block_01", "volumes": "2"},
          {"name": "group_app01", "status": "Error", "profile": "PF_block_02", "volumes": "5"}
      ];
      this.getVolumeGroups();
      this.volemeOptions = [
          {label: "group_for_REP",value:1},
          {label: "group_app01",value:2}
          ]
    }
    //show create volumes group
    createVolumeGroup(){
        this.volumeGroupForm.reset();
        this.showVolumeGroupDialog = true;
    }
    submit(group){
        if(!this.volumeGroupForm.valid){
            for(let i in this.volumeGroupForm.controls){
                this.volumeGroupForm.controls[i].markAsTouched();
            }
            return;
        }/*else{
            let param = {

            }
            this.volumeGroupService.createVolumeGroup(param).subscribe((res) => {
                this.getVolumeGroups();
            });
        }*/
        this.showVolumeGroupDialog = false;
    }
    getVolumeGroups(){
        this.volumeGroupService.getVolumeGroups().subscribe((res) => {
            console.log(res);
            this.volumeGroups = res.json();
        });
    }

}
