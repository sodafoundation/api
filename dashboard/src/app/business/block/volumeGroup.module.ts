import { NgModule, APP_INITIALIZER } from '@angular/core';
import { VolumeGroupComponent } from './volumeGroup.component';
import { ButtonModule, DataTableModule, InputTextModule } from '../../components/common/api';

import { HttpService } from './../../shared/service/Http.service';
import { VolumeService } from './volume.service';

@NgModule({
  declarations: [ VolumeGroupComponent ],
  imports: [ ButtonModule, DataTableModule, InputTextModule ],
  exports: [ VolumeGroupComponent ],
  providers: [
    HttpService,
    VolumeService
  ]
})
export class VolumeGroupModule { }