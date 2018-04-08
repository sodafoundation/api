import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ResourceComponent } from './resource.component';
import { RouterModule } from '@angular/router';
import { TabViewModule, ButtonModule } from '../../components/common/api';
import { VolumeListModule } from './volumeList.module';
import { VolumeGroupModule } from './volumeGroup.module';

let routers = [{
  path: '',
  component: ResourceComponent
}]

@NgModule({
  declarations: [
    ResourceComponent
  ],
  imports: [
    RouterModule.forChild(routers),
    VolumeListModule,
    VolumeGroupModule,
    TabViewModule,
    ButtonModule
  ],
  providers: []
})
export class ResourceModule { }