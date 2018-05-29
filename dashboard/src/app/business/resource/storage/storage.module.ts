import { NgModule, APP_INITIALIZER } from '@angular/core';
import { StorageComponent } from './storage.component';
import { ButtonModule, DataTableModule, InputTextModule } from './../../../components/common/api';


@NgModule({
  declarations: [
    StorageComponent
  ],
  imports: [
    ButtonModule,
    DataTableModule,
    InputTextModule
  ],
  exports: [
    StorageComponent
  ],
  providers: []
})
export class StorageModule { }