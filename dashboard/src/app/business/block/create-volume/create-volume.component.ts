import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { Validators, FormControl, FormGroup, FormBuilder } from '@angular/forms';

import { Message, SelectItem } from './../../../components/common/api';

import { VolumeService ,ReplicationService} from './../volume.service';
import { ProfileService } from './../../profile/profile.service';
import { AvailabilityZonesService } from './../../resource/resource.service';
import { I18NService,Utils } from 'app/shared/api';

@Component({
  selector: 'app-create-volume',
  templateUrl: './create-volume.component.html',
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
export class CreateVolumeComponent implements OnInit {

  label = {
    zone: this.i18n.keyID["sds_block_volume_az"],
    name: this.i18n.keyID["sds_block_volume_name"],
    profile: this.i18n.keyID["sds_block_volume_profile"],
    capacity: this.i18n.keyID["sds_home_capacity"],
    quantity: this.i18n.keyID["sds_block_volume_quantity"]
  };
  availabilityZones = [];
  volumeform;
  volumeItems = [0];
  capacityUnit = [];
  profileOptions = [];
  capacity = 'GB';
  createVolumes = [];
  value: boolean;
  showReplicationConf = false;
  errorMessage = {
    "zone": { required: "Zone is required."}
  };
    validRule= {
        'name':'^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){0,127}$'
    };
    defaultProfile = {
        label: null,
        value: {id:null,profileName:null}
    };
  constructor(
    private router: Router,
    private fb: FormBuilder,
    private ProfileService: ProfileService,
    private VolumeService: VolumeService,
    private replicationService:ReplicationService,
    private availabilityZonesService:AvailabilityZonesService,
    public i18n:I18NService
  ) {}

  ngOnInit() {
    this.getAZ();
    this.getProfiles();

    this.capacityUnit = [
      {
        label: 'GB', value: 'GB'
      },
      {
        label: 'TB', value: 'TB'
      }
    ];
    this.volumeform = this.fb.group({
      'zone': new FormControl('', Validators.required),
      'name0': new FormControl('', {validators:[Validators.required,Validators.pattern(this.validRule.name)]}),
      'profileId0': new FormControl(this.defaultProfile, {validators:[Validators.required,this.checkProfile]}),
      'size0': new FormControl(1, Validators.required),
      'capacity0': new FormControl(''),
      'quantity0': new FormControl(1)
    });
    this.volumeform.valueChanges.subscribe(
      (value:string)=>{
          this.createVolumes = this.getVolumesDataArray(this.volumeform.value);
          this.setRepForm();
      }
    );
      this.createVolumes = this.getVolumesDataArray(this.volumeform.value);
      this.setRepForm();
  }

  addVolumeItem() {
    this.volumeItems.push(
      this.volumeItems[this.volumeItems.length-1] + 1
    );
    this.volumeItems.forEach(index => {
      if(index !== 0){
        this.volumeform.addControl('name'+index, this.fb.control('', Validators.required));
        this.volumeform.addControl('profileId'+index, this.fb.control(this.defaultProfile,Validators.required));
        this.volumeform.addControl('size'+index, this.fb.control(1, Validators.required));
        this.volumeform.addControl('capacity'+index, this.fb.control('GB', Validators.required));
        this.volumeform.addControl('quantity'+index, this.fb.control(1));
      }
    });
  }
  getAZ(){
    this.availabilityZonesService.getAZ().subscribe((azRes) => {
      let AZs=azRes.json();
      let azArr = [];
      if(AZs && AZs.length !== 0){
          AZs.forEach(item =>{
              let obj = {label: item, value: item};
              azArr.push(obj);
          })
      }
      this.availabilityZones = azArr;
    })
  }
  getProfiles() {
    this.ProfileService.getProfiles().subscribe((res) => {
      let profiles = res.json();
      profiles.forEach(profile => {
        this.profileOptions.push({
          label: profile.name,
          value: {id:profile.id,profileName:profile.name}
        });
      });
    });
  }

  deleteVolumeItem(index) {
      this.volumeItems.splice(index, 1);
      this.volumeform.removeControl('name'+index);
      this.volumeform.removeControl('profileId'+index);
      this.volumeform.removeControl('size'+index);
      this.volumeform.removeControl('capacity'+index);
      this.volumeform.removeControl('quantity'+index);
  }

  createVolume(param){
    this.VolumeService.createVolume(param).subscribe((res) => {
      this.router.navigate(['/block']);
    });
  }
  createVolumeAndReplication(volParam,repParam){
    this.VolumeService.createVolume(volParam).subscribe((res2) => {
        this.VolumeService.createVolume(repParam).subscribe((res) => {
            let param = {
                "name":res.json().name ,
                "primaryVolumeId": res2.json().id,
                "availabilityZone": res.json().availabilityZone,
                "profileId": res.json().profileId,
                "replicationMode":"async",
                "replicationPeriod":this.createVolumes["formGroup"].value.period,
                "secondaryVolumeId":res.json().id
            }
            this.replicationService.createReplication(param).subscribe((res) => {});
            this.router.navigate(['/block']);
        });
    });
  }
  onSubmit(value) {
      if(!this.volumeform.valid){
          for(let i in this.volumeform.controls){
              this.volumeform.controls[i].markAsTouched();
          }
          return;
      }
      if(this.showReplicationConf && !this.createVolumes["formGroup"].valid){
          for(let i in this.createVolumes["formGroup"].controls){
              this.createVolumes["formGroup"].controls[i].markAsTouched();
          }
          return;
      }
      let dataArr = this.getVolumesDataArray(value);
      let volumeData = [];
      dataArr.forEach(item => {
          volumeData.push({
              name: item.name,
              size: item.size,
              availabilityZone: item.availabilityZone,
              profileId: item.profile.id
          });
      });
      for(let i in volumeData){
          if(this.showReplicationConf){
              let repVolume = {
                  name:null,
                  profileId:null,
                  availabilityZone: null
              };
              Object.assign(repVolume,volumeData[i]);
              repVolume.name = this.createVolumes["formGroup"].value["name"+i];
              repVolume.profileId = this.createVolumes["formGroup"].value["profileId"+i];
              repVolume.availabilityZone = "secondary";
              this.createVolumeAndReplication(volumeData[i],repVolume);
          }else{
              this.createVolume(volumeData[i]);
          }
      }
  }
  getVolumesDataArray(value){
      let dataArr = [];
      this.volumeItems.forEach(index => {
          if(!value['capacity'+index]){
              value['capacity'+index]='GB';
          }
          let unit = value['capacity'+index]==='GB' ? 1 : 1024;
          let qunantity = value['quantity'+index];
          if(qunantity && qunantity !== 1){
              for(let i=0;i<qunantity;i++){
                  dataArr.push({
                      name: value['name'+index]+i,
                      size: value['size'+index]*unit,
                      availabilityZone: value.zone,
                      profile: value['profileId'+index]
                  });
              }
          }else{
              dataArr.push({
                  name: value['name'+index],
                  size: value['size'+index]*unit,
                  availabilityZone: value.zone,
                  profile: value['profileId'+index]
              });
          }
      });
      return dataArr;
  }
    checkRep(param:boolean){}
    //create replication volumes formGroup
    setRepForm(){
        let param = {
            'zone': new FormControl(this.createVolumes[0].availabilityZone, Validators.required),
            'period': new FormControl(60, Validators.required)
        };
        for(let i in this.createVolumes){
            param["name"+i] = new FormControl(this.createVolumes[i].name+"-replication", Validators.required);
            param["profileId"+i] = new FormControl('', Validators.required);
        }
        this.createVolumes["formGroup"] = this.fb.group(param);
    }
    checkProfile(control:FormControl):{[s:string]:boolean}{
      if(control.value.id == null){
          return {profileNull:true}
      }
    }
    getErrorMessage(control,extraParam){
        let page = "";
        let key = Utils.getErrorKey(control,page);
        return extraParam ? this.i18n.keyID[key].replace("{0}",extraParam):this.i18n.keyID[key];
    }
}
