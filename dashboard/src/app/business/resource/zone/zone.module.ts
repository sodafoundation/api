import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ZoneComponent } from './zone.component';
import { ButtonModule,DataTableModule } from './../../../components/common/api';

@NgModule({
  declarations: [
    ZoneComponent
  ],
  imports: [
    ButtonModule,
    DataTableModule
  ],
  exports: [
    ZoneComponent
  ],
  providers: []
})
export class ZoneModule { }