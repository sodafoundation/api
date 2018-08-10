import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ZoneComponent } from './zone.component';
import { ButtonModule, DataTableModule, InputTextModule } from './../../../components/common/api';
import {AvailabilityZonesService} from '../resource.service';
import { HttpService } from '../../../shared/service/Http.service';

@NgModule({
  declarations: [
    ZoneComponent
  ],
  imports: [
    ButtonModule,
    DataTableModule,
    InputTextModule
  ],
  exports: [
    ZoneComponent
  ],
  providers: [HttpService,AvailabilityZonesService]
})
export class ZoneModule { }