import { Component, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CreateProfileComponent } from './createProfile.component';
import { RouterModule } from '@angular/router';
import { InputTextModule, CheckboxModule, FormModule, ButtonModule } from '../../../components/common/api';

let routers = [{
  path: '',
  component: CreateProfileComponent
}]

@NgModule({
  declarations: [
    CreateProfileComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(routers),
    InputTextModule,
    CheckboxModule,
    ButtonModule,
    // FormModule
  ],
  providers: []
})
export class CreateProfileModule { }