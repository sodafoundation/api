import { NgModule, APP_INITIALIZER } from '@angular/core';
import { VolumeListComponent } from './volumeList.component';
import { ButtonModule, DataTableModule, DropMenuModule } from '../../components/common/api';

import { HttpService } from './../../shared/service/Http.service';
import { VolumeService } from './volume.service';

@NgModule({
  declarations: [ VolumeListComponent ],
  imports: [ ButtonModule, DataTableModule, DropMenuModule ],
  exports: [ VolumeListComponent ],
  providers: [
    HttpService,
    VolumeService
  ]
})
export class VolumeListModule { }