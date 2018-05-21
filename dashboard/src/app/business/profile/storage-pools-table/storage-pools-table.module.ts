import { CommonModule } from '@angular/common';
import { NgModule, APP_INITIALIZER } from '@angular/core';

import { TableModule } from './../../../components/common/api';
import { HttpService } from './../../../shared/api';
import { PoolService } from './../profile.service';

import { StoragePoolsTableComponent } from './storage-pools-table.component';

@NgModule({
  declarations: [
    StoragePoolsTableComponent
  ],
  imports: [
    CommonModule,
    TableModule
  ],
  providers: [
    HttpService,
    PoolService
  ],
  exports: [
    StoragePoolsTableComponent
  ],
})
export class PoolModule { }