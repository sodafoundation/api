import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { VolumeListComponent } from './volumeList.component';
import { ButtonModule, DataTableModule, DropMenuModule, DialogModule, FormModule } from '../../components/common/api';

import { HttpService } from './../../shared/service/Http.service';
import { VolumeService } from './volume.service';

@NgModule({
  declarations: [ VolumeListComponent ],
  imports: [ 
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    ButtonModule,
    DataTableModule,
    DropMenuModule,
    DialogModule,
    FormModule
  ],
  exports: [ VolumeListComponent ],
  providers: [
    HttpService,
    VolumeService
  ]
})
export class VolumeListModule { }