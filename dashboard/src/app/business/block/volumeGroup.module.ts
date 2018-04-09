import { NgModule, APP_INITIALIZER } from '@angular/core';
import { VolumeGroupComponent } from './volumeGroup.component';
import { ButtonModule, DataTableModule } from '../../components/common/api';

@NgModule({
  declarations: [ VolumeGroupComponent ],
  imports: [ ButtonModule, DataTableModule ],
  exports: [ VolumeGroupComponent ],
  providers: []
})
export class VolumeGroupModule { }