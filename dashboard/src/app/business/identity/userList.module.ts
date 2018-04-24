import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { UserListComponent } from './userList.component';
import { ButtonModule, DataTableModule, DropMenuModule, DialogModule, InputTextModule, InputTextareaModule, DropdownModule, PasswordModule, ConfirmDialogModule } from '../../components/common/api';
import { UserDetailModule } from './userDetail/userDetail.module';

@NgModule({
  declarations: [UserListComponent],
  imports: [
    CommonModule,
    ButtonModule,
    DataTableModule,
    DropMenuModule,
    DialogModule,
    ConfirmDialogModule,
    InputTextModule,
    InputTextareaModule,  
    DropdownModule,
    PasswordModule,
    FormsModule,
    UserDetailModule
  ],
  exports: [UserListComponent],
  providers: []
})
export class UserListModule { }