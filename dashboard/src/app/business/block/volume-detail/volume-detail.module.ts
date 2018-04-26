import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { VolumeDetailComponent } from './volume-detail.component';

import { TabViewModule } from './../../../components/common/api';
import { HttpService } from './../../../shared/service/Http.service';
import { VolumeService } from './../volume.service';

let routers = [{
  path: '',
  component: VolumeDetailComponent
}]

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild(routers),
    TabViewModule
  ],
  declarations: [
    VolumeDetailComponent
  ],
  providers: [
    HttpService,
    VolumeService
  ]
})
export class VolumeDetailModule { }
