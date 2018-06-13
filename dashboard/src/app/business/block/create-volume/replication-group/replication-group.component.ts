import { Component, OnInit,Input,Output,EventEmitter} from '@angular/core';
import { trigger, state, style, transition, animate } from '@angular/animations';

import { ProfileService } from './../../../profile/profile.service';
import {FormBuilder, FormControl, Validators} from "@angular/forms";

@Component({
  selector: 'app-replication-group',
  templateUrl: './replication-group.component.html',
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
export class ReplicationGroupComponent implements OnInit {
  groupOptions = [];
  createOrExist = 'createGroup';
  existingGroupLists = [];
  replicationEnable:boolean;
  @Input() volumes:any;
  @Output() checkParam = new EventEmitter();
  period = 60;
  selectedProfile;
  availabilityZone=[];
  profileOptions = [
    {
      label: 'Select Profile',
      value: null
    }
  ];

  constructor(
    private ProfileService: ProfileService,
    private fb: FormBuilder
  ) { }

  ngOnInit() {
    this.groupOptions = [
      {
        label: 'Create Group',
        value: 'createGroup'
      },
      {
        label: 'Add to Existing Group',
        value: 'existingGroup'
      }
    ];
      this.availabilityZone = [
          {
              label: 'Secondary', value: 'secondary'
          }
      ];
    this.getProfiles();
  }

  getReplicationGroup(){
    if(this.createOrExist==='existingGroup'){
      this.existingGroupLists = [
        {
          label: 'group_for_REP',
          value: 'group_for_REP_id1'
        },
        {
          label: 'group_for_REP',
          value: 'group_for_REP_id2'
        }
      ]
    }
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
    checkParamSubmit(){
      this.checkParam.emit(false);
    }
}
