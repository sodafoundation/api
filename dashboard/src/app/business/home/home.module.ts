import { NgModule, APP_INITIALIZER } from '@angular/core';
import { HomeComponent } from './home.component';
import { RouterModule } from '@angular/router';
import { ButtonModule } from '../../components/common/api';

let routers = [{
  path: '',
  component: HomeComponent
}]

@NgModule({
  declarations: [
    HomeComponent
  ],
  imports: [ RouterModule.forChild(routers), ButtonModule],
  providers: []
})
export class HomeModule { }