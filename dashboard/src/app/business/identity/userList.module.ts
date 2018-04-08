import { NgModule, APP_INITIALIZER } from '@angular/core';
import { UserListComponent } from './userList.component';
import { ButtonModule, DataTableModule, DropMenuModule } from '../../components/common/api';

@NgModule({
  declarations: [ UserListComponent ],
  imports: [ ButtonModule, DataTableModule, DropMenuModule ],
  exports: [ UserListComponent ],
  providers: []
})
export class UserListModule { }