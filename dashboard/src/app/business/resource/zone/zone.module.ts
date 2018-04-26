import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ZoneComponent } from './zone.component';
import { ButtonModule, DataTableModule, InputTextModule } from './../../../components/common/api';

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
  providers: []
})
export class ZoneModule { }