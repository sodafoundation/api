import { Component, OnInit } from '@angular/core';
import { Router,ActivatedRoute} from '@angular/router';

import { VolumeService } from './../volume.service';
import { ProfileService } from './../../profile/profile.service';

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

  constructor(
    private VolumeService: VolumeService,
    private ActivatedRoute: ActivatedRoute,
    private ProfileService: ProfileService
  ) { }

  ngOnInit() {
    this.ActivatedRoute.params.subscribe((params) => this.volumeId = params.volumeId);

    this.items = [
      { label: 'Volume', url: '/block' },
      { label: 'Volume detail', url: '/volumeDetail' }
    ];

    this.label = {
      Name: 'Name',
      Profile: 'Profile',
      Status: 'Status',
      VolumeID: 'Volume ID',
      Capacity: 'Capacity',
      CreatedAt: 'Created At'
    };

    this.getVolume(this.volumeId);
  }

  getVolume(id){
    this.VolumeService.getVolumeById(id).subscribe((res) => {
      this.volume = res.json();
      this.ProfileService.getProfileById(this.volume.profileId).subscribe((res)=>{
          this.volume.profileName = res.json().name;
      })
    });
  }

}
