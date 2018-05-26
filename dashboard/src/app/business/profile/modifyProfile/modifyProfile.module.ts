import { Component, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { modifyProfileComponent } from './modifyProfile.component';
import { RouterModule } from '@angular/router';
import { BreadcrumbModule,ChartModule,ButtonModule } from './../../../components/common/api';

import { HttpService } from './../../../shared/api';
import { ProfileService,PoolService } from './../profile.service';
import { PoolModule } from './../storage-pools-table/storage-pools-table.module';

let routers = [{
  path: '',
  component: modifyProfileComponent
}]

@NgModule({
  declarations: [
    modifyProfileComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(routers),
    BreadcrumbModule,
    ChartModule,
    ButtonModule,
    PoolModule
    // FormModule
  ],
  providers: [
    HttpService,
    ProfileService,
    PoolService
  ]
})
export class ModifyProfileModule { }