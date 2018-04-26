import { Component, OnInit } from '@angular/core';
import { Router,ActivatedRoute} from '@angular/router';

import { VolumeService } from './../volume.service';

@Component({
  selector: 'app-volume-detail',
  templateUrl: './volume-detail.component.html',
  styleUrls: ['./volume-detail.component.scss']
})
export class VolumeDetailComponent implements OnInit {
  items;
  label;
  volume;
  volumeId;

  constructor(
    private VolumeService: VolumeService,
    private ActivatedRoute: ActivatedRoute
  ) { }

  ngOnInit() {
    let volumeId;
    this.ActivatedRoute.params.subscribe((params) => volumeId = params.volumeId);

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

    this.getVolume(volumeId);
  }

  getVolume(id){
    this.VolumeService.getVolumeById(id).subscribe((res) => {
      this.volume = res.json();
      console.log(this.volume);
    });
  }

}
