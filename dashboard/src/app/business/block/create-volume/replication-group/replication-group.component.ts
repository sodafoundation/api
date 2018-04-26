import { Component, OnInit } from '@angular/core';

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

  constructor() { }

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

}
