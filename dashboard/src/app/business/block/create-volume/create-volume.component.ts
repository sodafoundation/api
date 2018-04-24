import { Component, OnInit } from '@angular/core';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { Validators, FormControl, FormGroup, FormBuilder } from '@angular/forms';

import { Message, SelectItem } from './../../../components/common/api';

@Component({
  selector: 'app-create-volume',
  templateUrl: './create-volume.component.html',
  styleUrls: ['./create-volume.component.scss'],
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

  label = {};
  availabilityZones = [];
  volumeform;
  volumeItems = [{}];
  capacityUnit = [];

  value: boolean;

  constructor(
    private fb: FormBuilder
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
        label: 'Select Zone', value: null
      },
      {
        label: 'DataCenter_UnitA', value: '1'
      },
      {
        label: 'DataCenter_UnitB', value: '2'
      }
    ];

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
      'name': new FormControl('', Validators.required),
      'capacity': new FormControl('', Validators.required),
      'quantity': new FormControl('')
    });
  }

  addVolumeItem(){
    this.volumeItems.push({});
  }

  deleteVolumeItem(index){
    this.volumeItems.splice(index,1);
  }

}
