import { NgModule, APP_INITIALIZER } from '@angular/core';
import { VolumeListComponent } from './volumeList.component';
import { ButtonModule, DataTableModule, DropMenuModule } from '../../components/common/api';

@NgModule({
  declarations: [ VolumeListComponent ],
  imports: [ ButtonModule, DataTableModule, DropMenuModule ],
  exports: [ VolumeListComponent ],
  providers: []
})
export class VolumeListModule { }