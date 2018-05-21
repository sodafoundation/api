import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';

import { ButtonModule,MessageModule,TabViewModule,DialogModule,DropdownModule } from '../../components/common/api';

import { CloudComponent } from './cloud.component';
import { RegistryComponent } from './registry/registry.component';
import { ReplicationComponent } from './replication/replication.component';
import { MigrationComponent } from './migration/migration.component';
import { CloudServiceItemComponent } from './cloud-service-item/cloud-service-item.component';

let routers = [{
  path: '',
  component: CloudComponent
}]

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild(routers),
    FormsModule,
    ButtonModule,
    MessageModule,
    TabViewModule,
    DialogModule,
    DropdownModule
  ],
  declarations: [
    CloudComponent,
    RegistryComponent,
    ReplicationComponent,
    MigrationComponent,
    CloudServiceItemComponent
  ]
})
export class CloudModule { }
