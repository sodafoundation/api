import { NgModule, APP_INITIALIZER } from '@angular/core';
import { RegionComponent } from './region.component';
import { ButtonModule, DataTableModule, InputTextModule } from './../../../components/common/api';


@NgModule({
  declarations: [
    RegionComponent
  ],
  imports: [
    ButtonModule,
    DataTableModule,
    InputTextModule
  ],
  exports: [ RegionComponent ],
  providers: []
})
export class RegionModule { }