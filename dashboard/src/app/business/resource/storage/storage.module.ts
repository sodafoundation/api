import { NgModule, APP_INITIALIZER } from '@angular/core';
import { StorageComponent } from './storage.component';
import { ButtonModule,DataTableModule } from './../../../components/common/api';


@NgModule({
  declarations: [
    StorageComponent
  ],
  imports: [
    ButtonModule,
    DataTableModule
  ],
  exports: [
    StorageComponent
  ],
  providers: []
})
export class StorageModule { }