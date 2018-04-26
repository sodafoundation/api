import { NgModule, APP_INITIALIZER } from '@angular/core';
import { RegionComponent } from './region.component';
import { ButtonModule,DataTableModule } from './../../../components/common/api';


@NgModule({
  declarations: [
    RegionComponent
  ],
  imports: [
    ButtonModule,
    DataTableModule
  ],
  exports: [ RegionComponent ],
  providers: []
})
export class RegionModule { }