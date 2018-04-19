import { NgModule, APP_INITIALIZER } from '@angular/core';
import { HomeComponent } from './home.component';
import { ImgItemComponent } from './imgItem.component/imgItem.component';

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
    ChartModule
  ],
  providers: []
})
export class HomeModule { }