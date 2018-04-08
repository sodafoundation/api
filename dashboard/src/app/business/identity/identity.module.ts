import { NgModule, APP_INITIALIZER } from '@angular/core';
import { IdentityComponent } from './identity.component';
import { RouterModule } from '@angular/router';
import { TabViewModule, ButtonModule, DialogModule} from '../../components/common/api';
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
    TabViewModule,
    ButtonModule,
    DialogModule
  ],
  providers: []
})
export class IdentityModule { }