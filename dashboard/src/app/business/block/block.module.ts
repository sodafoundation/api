import { NgModule, APP_INITIALIZER } from '@angular/core';
import { BlockComponent } from './block.component';
import { RouterModule } from '@angular/router';
import { TabViewModule, ButtonModule } from '../../components/common/api';
import { VolumeListModule } from './volumeList.module';
import { VolumeGroupModule } from './volumeGroup.module';
import { CreateVolumeGroupComponent } from './create-volume-group/create-volume-group.component';

let routers = [{
  path: '',
  component: BlockComponent
}]

@NgModule({
  declarations: [
    BlockComponent,
    CreateVolumeGroupComponent
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
export class BlockModule { }