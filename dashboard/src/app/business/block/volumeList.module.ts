import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { VolumeListComponent } from './volumeList.component';
import { ButtonModule, DataTableModule, DropMenuModule, DialogModule, FormModule, InputTextModule, InputTextareaModule, DropdownModule } from '../../components/common/api';

import { HttpService } from './../../shared/service/Http.service';
import { VolumeService,SnapshotService } from './volume.service';
import { ProfileService } from './../profile/profile.service';

@NgModule({
  declarations: [ VolumeListComponent ],
  imports: [ 
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    ButtonModule,
    InputTextModule,
    InputTextareaModule,
    DataTableModule,
    DropdownModule,
    DropMenuModule,
    DialogModule,
    FormModule
  ],
  exports: [ VolumeListComponent ],
  providers: [
    HttpService,
    VolumeService,
    SnapshotService,
    ProfileService
  ]
})
export class VolumeListModule { }