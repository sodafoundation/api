import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HomeComponent } from './home.component';
import { ImgItemComponent } from './imgItem.component/imgItem.component';
import { ProfileService } from 'app/business/profile/profile.service';
import { HttpService } from 'app/shared/service/Http.service';

import { RouterModule } from '@angular/router';
import { ButtonModule, ChartModule } from '../../components/common/api';

let routers = [{
  path: '',
  component: HomeComponent
}]

@NgModule({
  declarations: [
    HomeComponent,
    ImgItemComponent,
  ],
  imports: [
    RouterModule.forChild(routers), ButtonModule,
    CommonModule,
    ChartModule
  ],
  providers: [HttpService, ProfileService]
})
export class HomeModule { }