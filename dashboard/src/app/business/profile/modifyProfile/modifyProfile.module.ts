import { Component, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { modifyProfileComponent } from './modifyProfile.component';
import { RouterModule } from '@angular/router';
import { BreadcrumbModule,ChartModule,TableModule,ButtonModule } from './../../../components/common/api';

let routers = [{
  path: '',
  component: modifyProfileComponent
}]

@NgModule({
  declarations: [
    modifyProfileComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(routers),
    BreadcrumbModule,
    ChartModule,
    TableModule,
    ButtonModule
    // FormModule
  ],
  providers: []
})
export class ModifyProfileModule { }