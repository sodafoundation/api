import { Component, OnInit } from '@angular/core';
import { Router,ActivatedRoute} from '@angular/router';

import { VolumeService } from './../volume.service';
import { ProfileService } from './../../profile/profile.service';
import { I18NService, Utils } from 'app/shared/api';

@Component({
  selector: 'app-volume-detail',
  templateUrl: './volume-detail.component.html',
  styleUrls: [

  ]
})
export class VolumeDetailComponent implements OnInit {
  items;
  label;
  volume;
  volumeId;
  showVolumeSource: boolean = false;
  volumeSource: string = "";

  constructor(
    private VolumeService: VolumeService,
    private ActivatedRoute: ActivatedRoute,
    private ProfileService: ProfileService,
    public i18n:I18NService
  ) { }

  ngOnInit() {
    this.ActivatedRoute.params.subscribe((params) => this.volumeId = params.volumeId);

    this.items = [
      { label: this.i18n.keyID["sds_block_volume_title"], url: '/block' },
      { label: this.i18n.keyID["sds_block_volume_detail"], url: '/volumeDetail' }
    ];

    this.label = {
      Name: this.i18n.keyID["sds_block_volume_name"],
      Profile: this.i18n.keyID["sds_block_volume_profile"],
      Status: this.i18n.keyID["sds_block_volume_status"],
      VolumeID: this.i18n.keyID["sds_block_volume_id"],
      Capacity: this.i18n.keyID["sds_home_capacity"],
      CreatedAt: this.i18n.keyID["sds_block_volume_createat"]
    };

    this.getVolume(this.volumeId);
  }

  getVolume(id){
    this.VolumeService.getVolumeById(id).subscribe((res) => {
      this.volume = res.json();
      this.volume.size = Utils.getDisplayGBCapacity(res.json().size);
      this.ProfileService.getProfileById(this.volume.profileId).subscribe((res)=>{
          this.volume.profileName = res.json().name;
      })

      if(this.volume.snapshotId != ""){
        this.showVolumeSource = true;
        this.volumeSource = this.i18n.keyID['sds_block_volume_source'].replace("{{}}", this.volume.snapshotId);
      }
    });
  }

}
