import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { Validators, FormControl, FormGroup, FormBuilder } from '@angular/forms';

import { Message, SelectItem } from './../../../components/common/api';

import { VolumeService } from './../volume.service';
import { ProfileService } from './../../profile/profile.service';

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

  bbbbb = 'name';
  label = {};
  availabilityZones = [];
  volumeform;
  volumeItems = [0];
  capacityUnit = [];
  profileOptions = [
    {
      label: 'Select Profile',
      value: null
    }
  ];
  capacity = 'GB';

  value: boolean;

  errorMessage = {
    "zone": { required: "Zone is required."},//已经默认了一个选项，不会出现这个错误
  };

  constructor(
    private router: Router,
    private fb: FormBuilder,
    private ProfileService: ProfileService,
    private VolumeService: VolumeService
  ) { }

  ngOnInit() {
    this.label = {
      zone: 'Availability Zone',
      name: 'Name',
      profile: 'Profile',
      capacity: 'Capacity',
      quantity: 'Quantity'
    }

    this.availabilityZones = [
      {
        label: 'Default', value: 'default'
      }
    ];

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
      'zone': new FormControl('default', Validators.required),
      'name0': new FormControl('', Validators.required),
      'profileId0': new FormControl('', Validators.required),
      'size0': new FormControl('', Validators.required),
      'capacity0': new FormControl(''),
      'quantity0': new FormControl('')
    });

  }

  addVolumeItem() {
    this.volumeItems.push(
      this.volumeItems[this.volumeItems.length-1] + 1
    );
    this.volumeItems.forEach(index => {
      if(index !== 0){
        this.volumeform.addControl('name'+index, this.fb.control('', Validators.required));
        this.volumeform.addControl('profileId'+index, this.fb.control('', Validators.required));
        this.volumeform.addControl('size'+index, this.fb.control('', Validators.required));
        this.volumeform.addControl('capacity'+index, this.fb.control('', Validators.required));
        this.volumeform.addControl('quantity'+index, this.fb.control(''));
      }
    });
  }

  getProfiles() {
    this.ProfileService.getProfiles().subscribe((res) => {
      let profiles = res.json();
      profiles.forEach(profile => {
        this.profileOptions.push({
          label: profile.name,
          value: profile.id
        });
      });
    });
  }

  deleteVolumeItem(index) {
    this.volumeItems.splice(index, 1);
  }

  createVolume(param){
    this.VolumeService.createVolume(param).subscribe((res) => {
      this.router.navigate(['/block']);
  });
  }

  onSubmit(value) {    
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
            profileId: value['profileId'+index]
          });
        }
      }else{
        dataArr.push({
          name: value['name'+index],
          size: value['size'+index]*unit,
          availabilityZone: value.zone,
          profileId: value['profileId'+index]
        });
      }
    });

    dataArr.forEach(data=>{
      this.createVolume(data);
    });
  }

}
