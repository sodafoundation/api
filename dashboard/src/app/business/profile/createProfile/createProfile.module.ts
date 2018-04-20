import { Component, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { CreateProfileComponent } from './createProfile.component';
import { RouterModule } from '@angular/router';
import { InputTextModule, CheckboxModule, FormModule, ButtonModule, DropdownModule, RadioButtonModule, DialogModule } from '../../../components/common/api';

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
    FormsModule,  
    RouterModule.forChild(routers),
    InputTextModule,
    CheckboxModule,
    ButtonModule,
    DropdownModule,
    RadioButtonModule,
    DialogModule
  ],
  providers: []
})
export class CreateProfileModule { }