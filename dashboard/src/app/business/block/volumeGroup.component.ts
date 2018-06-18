import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { trigger, state, style, transition, animate} from '@angular/animations';
import { DialogModule } from '../../components/common/api';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl} from '@angular/forms';
import { VolumeService ,VolumeGroupService} from './volume.service';
import { ProfileService } from './../profile/profile.service';
import { ConfirmationService,ConfirmDialogModule} from '../../components/common/api';
import { Router } from '@angular/router';

@Component({
    selector: 'volume-group-list',
    templateUrl: 'volumeGroup.html',
    providers: [ConfirmationService],
    styleUrls: [],
    animations: []
})
export class VolumeGroupComponent implements OnInit{
    volumeGroups=[];
    volemeOptions = [];
    profileOptions = [];
    currentGroup:any;
    availabilityZones = [];
    selectedOption :string;
    selectedVolumeGroups=[];
    profileJson = {};
    showVolumeGroupDialog :boolean = false;
    showModifyGroup :boolean = false;
    volumeGroupForm:any;
    modifyGroupForm:any;
    validRule= {
        'name':'^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){0,127}$'
    };
    constructor(
        public I18N: I18NService,
         private router: Router,
        private volumeGroupService : VolumeGroupService,
        private fb : FormBuilder,
        private profileService :ProfileService,
        private confirmationService:ConfirmationService
    ){
        this.volumeGroupForm = this.fb.group({
            "group_name":["",{validators:[Validators.required, Validators.pattern(this.validRule.name)], updateOn:'change'} ],
            "description":[""],
            "profile":["",Validators.required],
            "zone":[""]
        });
        this.modifyGroupForm = this.fb.group({
            "group_name":["",{validators:[Validators.required, Validators.pattern(this.validRule.name)], updateOn:'change'} ],
            "description":[""]
        });
    }
    errorMessage = {
        "group_name": { required: "group name is required.", pattern:"Beginning with a letter with a length of 1-128, it can contain letters / numbers / underlines."},
        "profile": { required: "profile is required."}
    };
    
    label = {
        group_name_lable:'Group Name',
        profile_label:'Profile',
        description:this.I18N.keyID['sds_block_volume_descri'],
        zone:this.I18N.keyID['sds_block_volume_az']
    }

    ngOnInit() {
      this.availabilityZones = [
        {
          label: 'Default', value: 'default'
        }
      ];
      this.volemeOptions = [];
      this.getProfiles();
    }
    //show create volumes group
    createVolumeGroup(){
        this.volumeGroupForm.reset();
        this.showVolumeGroupDialog = true;
    }
    ModifyVolumeGroupDisplay(volumeGroup){
        this.modifyGroupForm.reset();
        this.currentGroup = volumeGroup;
        this.modifyGroupForm.controls['group_name'].setValue(this.currentGroup.name);
        this.modifyGroupForm.controls['description'].setValue("");
        this.showModifyGroup = true;
    }
    submit(group){
        if(!this.volumeGroupForm.valid){
            for(let i in this.volumeGroupForm.controls){
                this.volumeGroupForm.controls[i].markAsTouched();
            }
            return;
        }else{
            let param = {
                name : group.group_name,
                profiles : group.profile,
                description:group.description,
                availabilityZone:group.zone
            }
            this.volumeGroupService.createVolumeGroup(param).subscribe((res) => {
                this.getVolumeGroups();
            });
        }
        this.showVolumeGroupDialog = false;
    }
    getVolumeGroups(){
        this.volumeGroupService.getVolumeGroups().subscribe((res) => {
            let volumeGroups = res.json();
            if(volumeGroups && volumeGroups.length != 0){
                volumeGroups.forEach((item)=>{
                    if(!item.description){
                        item.description = "--";
                    }
                    let profileName = [];
                    item.profiles.forEach((profileId)=>{
                        profileName.push(this.profileJson[profileId]);
                    });
                    item.profileName = profileName;
                });
            }
            this.volumeGroups = volumeGroups;
        });
    }
    getProfiles() {
        this.profileService.getProfiles().subscribe((res) => {
            let profiles = res.json();
            profiles.forEach(profile => {
                this.profileOptions.push({
                    label: profile.name,
                    value: profile.id
                });
                this.profileJson[profile.id] = profile.name;
            });
            this.getVolumeGroups();
        });
    }
    deleteVolumeGroup(volumeGroup){
        this.currentGroup = volumeGroup;
        let msg = "<div>Are you sure you want to delete the Volume Group?</div><h3>[ "+ volumeGroup.name +" ]</h3>";
        let header ="Delete Volume Group";
        let acceptLabel = "Delete";
        let warming = true;
        this.confirmDialog([msg,header,acceptLabel,warming,"delete"])
    }
    deleteMultiVolumeGroups(){
        let msg = "<div>Are you sure you want to delete the selected Volume Groups?</div><h3>[ "+ this.selectedVolumeGroups.length +" Volume Group ]</h3>";
        let header ="Delete Volume Group";
        let acceptLabel = "Delete";
        let warming = true;
        this.confirmDialog([msg,header,acceptLabel,warming,"multiDelete"])
    }
    modifyGroup(value){
        if(!this.modifyGroupForm.valid){
            for(let i in this.modifyGroupForm.controls){
                this.modifyGroupForm.controls[i].markAsTouched();
            }
            return;
        }else{
            let param = {
                name:value.group_name,
                description:value.description,
            }
            this.volumeGroupService.modifyVolumeGroup(this.currentGroup.id,param).subscribe((res) => {
                this.getVolumeGroups();
            });
        }
        this.showModifyGroup = false;
    }
    confirmDialog([msg,header,acceptLabel,warming=true,func]){
        this.confirmationService.confirm({
            message: msg,
            header: header,
            acceptLabel: acceptLabel,
            isWarning: warming,
            accept: ()=>{
                try {
                    if(func === "delete"){
                        this.volumeGroupService.deleteVolumeGroup(this.currentGroup.id).subscribe((res) => {
                            this.getVolumeGroups();
                        })
                    }else if(func === "multiDelete"){
                        this.selectedVolumeGroups.forEach(item=>{
                            this.volumeGroupService.deleteVolumeGroup(item.id).subscribe((res) => {
                                this.getVolumeGroups();
                            })
                        });
                        this.selectedVolumeGroups = [];
                    }
                }
                catch (e) {
                    console.log(e);
                }
                finally {
                    
                }
            },
            reject:()=>{}
        })
    }

}
