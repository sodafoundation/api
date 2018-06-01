import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { userDetailComponent } from './userDetail.component';
import { ButtonModule, DataTableModule, DropMenuModule } from '../../../components/common/api';

@NgModule({
  declarations: [ userDetailComponent ],
  imports: [ CommonModule, ButtonModule, DataTableModule, DropMenuModule ],
  exports: [ userDetailComponent ],
  providers: []
})
export class UserDetailModule { }