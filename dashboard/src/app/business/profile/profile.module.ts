import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ProfileComponent } from './profile.component';
import { RouterModule } from '@angular/router';
import { ButtonModule } from '../../components/common/api';

let routers = [{
  path: '',
  component: ProfileComponent
}]

@NgModule({
  declarations: [
    ProfileComponent
  ],
  imports: [ RouterModule.forChild(routers), ButtonModule ],
  providers: []
})
export class ProfileModule { }