import { Component, OnInit } from '@angular/core';
import { PoolService } from './../profile.service';

@Component({
  selector: 'app-storage-pools-table',
  templateUrl: './storage-pools-table.component.html',
  styleUrls: [
    
  ]
})
export class StoragePoolsTableComponent implements OnInit {
  cols;
  totalFreeCapacity = 0;
  pools = [];

  constructor(
    private PoolService: PoolService
  ) { }

  getPools() {
    this.PoolService.getPools().subscribe((res) => {
      this.pools = res.json();
      this.totalFreeCapacity = this.getSumCapacity(this.pools, 'free');
      console.log(res.json());
    });
  }

  getSumCapacity(pools, FreeOrTotal) {
    let SumCapacity: number = 0;
    let arrLength = pools.length;
    for (let i = 0; i < arrLength; i++) {
      if (FreeOrTotal === 'free') {
        SumCapacity += pools[i].freeCapacity;
      } else {
        SumCapacity += pools[i].totalCapacity;
      }
    }
    return SumCapacity;
  }

  ngOnInit() {

    this.cols = [
      { field: 'name', header: 'Name' },
      { field: 'freeCapacity', header: 'FreeCapacity' },
      { field: 'totalCapacity', header: 'TotalCapacity' },
      { field: 'extras.advanced.diskType', header: 'Disk' },
      { field: 'extras.dataStorage.provisioningPolicy', header: 'Backend' }
    ];
    this.getPools();
  }

}
