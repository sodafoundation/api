import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ResourceComponent } from './resource.component';
import { RouterModule } from '@angular/router';
import { TabViewModule, ButtonModule } from '../../components/common/api';

// 引入region模块
import { RegionModule } from './region/region.module';
// 引入zone模块
import { ZoneModule } from './zone/zone.module';
// 引入storage模块
import { StorageModule } from './storage/storage.module';

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
    TabViewModule,
    ButtonModule,
    RegionModule,//region模块
    ZoneModule,//zone模块
    StorageModule//storage模块
  ],
  providers: []
})
export class ResourceModule { }