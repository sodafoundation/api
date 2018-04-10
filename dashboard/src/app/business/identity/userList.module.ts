import { NgModule, APP_INITIALIZER } from '@angular/core';
import { UserListComponent } from './userList.component';
import { ButtonModule, DataTableModule, DropMenuModule, DialogModule, InputTextModule, InputTextareaModule, DropdownModule, PasswordModule } from '../../components/common/api';

@NgModule({
  declarations: [UserListComponent],
  imports: [
    ButtonModule,
    DataTableModule,
    DropMenuModule,
    DialogModule,
    InputTextModule,
    InputTextareaModule,  
    DropdownModule,
    PasswordModule
  ],
  exports: [UserListComponent],
  providers: []
})
export class UserListModule { }