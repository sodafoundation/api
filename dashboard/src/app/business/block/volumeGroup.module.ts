import { NgModule, APP_INITIALIZER } from '@angular/core';
import { VolumeGroupComponent } from './volumeGroup.component';
import { ButtonModule, DataTableModule, InputTextModule, DialogModule,FormModule,MultiSelectModule ,DropdownModule} from '../../components/common/api';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { HttpService } from './../../shared/service/Http.service';
import { VolumeService ,VolumeGroupService} from './volume.service';

@NgModule({
  declarations: [ VolumeGroupComponent ],
  imports: [ ButtonModule, DataTableModule, InputTextModule, DialogModule,FormModule,MultiSelectModule,DropdownModule,ReactiveFormsModule,FormsModule],
  exports: [ VolumeGroupComponent ],
  providers: [
    HttpService,
    VolumeService,
    VolumeGroupService
  ]
})
export class VolumeGroupModule { }
