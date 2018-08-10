import { Component, OnInit,Input } from '@angular/core';
import { PoolService } from './../profile.service';
import { Utils } from '../../../shared/api';

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
  selectData;
  @Input()
  set selectProfile(selectProfile: any) {
    this.selectData = selectProfile;
    this.getPools();
  };
  constructor(
    private PoolService: PoolService
  ) { }

  getPools() {
    this.PoolService.getPools().subscribe((res) => {
        this.pools = [];
        let pools = res.json();
      if(this.selectData){
          let arrLength = pools.length
          for (let i = 0; i < arrLength; i++) {
              if (this.selectData.extras.protocol.toLowerCase() == pools[i].extras.ioConnectivity.accessProtocol &&  this.selectData.storageType == pools[i].extras.dataStorage.provisioningPolicy){
                this.pools.push(pools[i]);
              }
          }
      }else{
          this.pools = pools;
      }

      this.pools.map((pool)=>{
        pool.freeCapacityFormat = Utils.getDisplayGBCapacity(pool.freeCapacity);
        pool.totalCapacityFormat = Utils.getDisplayGBCapacity(pool.totalCapacity);
      })
      
      this.totalFreeCapacity = this.getSumCapacity(this.pools, 'free');
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
