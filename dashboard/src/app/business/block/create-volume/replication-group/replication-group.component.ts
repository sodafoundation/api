import { Component, OnInit } from '@angular/core';
import { ProfileService } from './../../../profile/profile.service';

@Component({
  selector: 'app-replication-group',
  templateUrl: './replication-group.component.html',
  styleUrls: ['./replication-group.component.scss']
})
export class ReplicationGroupComponent implements OnInit {
  groupOptions = [];
  createOrExist = 'createGroup';
  existingGroupLists = [];
  replicationEnable:boolean;
  volumes = [];
  period = 60;
  selectedProfile;
  profileOptions = [
    {
      label: 'Select Profile',
      value: null
    }
  ];

  constructor(
    private ProfileService: ProfileService
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
    this.volumes = [
      {
        availabilityZone: 'Region-Beijing',
        name: 'vol_01',
        size: 1,
        profileId: 'PF_block_01'
      },
      {
        availabilityZone: 'Region-Beijing',
        name: 'vol_02',
        size: 1,
        profileId: 'PF_block_02'
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

}
