import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { VolumeGroupDetailComponent } from './volume-group-detail.component';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { TabViewModule,ButtonModule, DataTableModule, DropMenuModule, DialogModule, FormModule, InputTextModule, InputTextareaModule, ConfirmDialogModule ,ConfirmationService} from './../../../components/common/api';
import { HttpService } from './../../../shared/service/Http.service';
import { VolumeService,VolumeGroupService} from './../volume.service';
import { ProfileService } from './../../profile/profile.service';

let routers = [{
  path: '',
  component: VolumeGroupDetailComponent
}]
@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild(routers),
    ReactiveFormsModule,
    FormsModule,
    InputTextModule,
    InputTextareaModule,
    TabViewModule,
    ButtonModule,
    DataTableModule,
    DialogModule,
    FormModule,
    ConfirmDialogModule
  ],
  declarations: [VolumeGroupDetailComponent],
  providers: [
    HttpService,
    VolumeService,
    ConfirmationService,
    ProfileService,
    VolumeGroupService
  ]
})
export class VolumeGroupDetailModule { }
