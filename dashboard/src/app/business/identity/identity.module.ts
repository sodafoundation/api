import { NgModule, APP_INITIALIZER } from '@angular/core';
import { IdentityComponent } from './identity.component';
import { RouterModule } from '@angular/router';
import { TabViewModule } from '../../components/common/api';
import { TenantListModule } from './tenantList.module';
import { UserListModule } from './userList.module';

let routers = [{
  path: '',
  component: IdentityComponent
}]

@NgModule({
  declarations: [
    IdentityComponent
  ],
  imports: [
    RouterModule.forChild(routers),
    TenantListModule,
    UserListModule,
    TabViewModule
  ],
  providers: []
})
export class IdentityModule { }