import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ServiceComponent } from './service.component';
import { RouterModule } from '@angular/router';

let routers = [{
  path: '',
  component: ServiceComponent
}]

@NgModule({
  declarations: [
    ServiceComponent
  ],
  imports: [ RouterModule.forChild(routers) ],
  providers: []
})
export class ServiceModule { }