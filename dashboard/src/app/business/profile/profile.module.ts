import { CommonModule } from '@angular/common';
import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ProfileComponent } from './profile.component';
import { RouterModule } from '@angular/router';

import { ProfileCardComponent } from './profileCard/profile-card.component';
import { ButtonModule,CardModule,ChartModule } from '../../components/common/api';

let routers = [{
  path: '',
  component: ProfileComponent
}]

@NgModule({
  declarations: [
    ProfileComponent,
    ProfileCardComponent
  ],
  imports: [
    RouterModule.forChild(routers),
    ButtonModule,
    CommonModule,
    CardModule,
    ChartModule
  ],
  providers: []
})
export class ProfileModule { }