import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { VolumeDetailComponent } from './volume-detail.component';

import { TabViewModule,ButtonModule, DataTableModule, DropMenuModule, DialogModule, FormModule, InputTextModule, InputTextareaModule, ConfirmDialogModule ,ConfirmationService} from './../../../components/common/api';
import { HttpService } from './../../../shared/service/Http.service';
import { VolumeService,SnapshotService ,ReplicationService} from './../volume.service';
import { SnapshotListComponent } from './snapshot-list/snapshot-list.component';
import { ReplicationListComponent } from './replication-list/replication-list.component';
import { ProfileService } from './../../profile/profile.service';

let routers = [{
  path: '',
  component: VolumeDetailComponent
}]

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    InputTextModule,
    InputTextareaModule,
    RouterModule.forChild(routers),
    TabViewModule,
    ButtonModule,
    DataTableModule,
    DialogModule,
    FormModule,
    ConfirmDialogModule
  ],
  declarations: [
    VolumeDetailComponent,
    SnapshotListComponent,
    ReplicationListComponent
  ],
  providers: [
    HttpService,
    VolumeService,
    SnapshotService,
    ConfirmationService,
    ProfileService,
    ReplicationService
  ]
})
export class VolumeDetailModule { }
