import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { VolumeDetailComponent } from './volume-detail.component';

import { TabViewModule,ButtonModule, DataTableModule, DropMenuModule, DialogModule, FormModule } from './../../../components/common/api';
import { HttpService } from './../../../shared/service/Http.service';
import { VolumeService,SnapshotService } from './../volume.service';
import { SnapshotListComponent } from './snapshot-list/snapshot-list.component';
import { ReplicationListComponent } from './replication-list/replication-list.component';

let routers = [{
  path: '',
  component: VolumeDetailComponent
}]

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    RouterModule.forChild(routers),
    TabViewModule,
    ButtonModule,
    DataTableModule,
    DialogModule,
    FormModule
  ],
  declarations: [
    VolumeDetailComponent,
    SnapshotListComponent,
    ReplicationListComponent
  ],
  providers: [
    HttpService,
    VolumeService,
    SnapshotService
  ]
})
export class VolumeDetailModule { }
