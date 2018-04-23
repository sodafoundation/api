import { Component, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { CreateProfileComponent } from './createProfile.component';
import { RouterModule } from '@angular/router';
import { InputTextModule, CheckboxModule, FormModule, ButtonModule, DropdownModule, RadioButtonModule, DialogModule, Message, GrowlModule, TableModule } from '../../../components/common/api';
import { HttpService } from './../../../shared/api';
import { ProfileService } from './../profile.service';

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
    ReactiveFormsModule,
    FormsModule,  
    RouterModule.forChild(routers),
    InputTextModule,
    CheckboxModule,
    ButtonModule,
    DropdownModule,
    RadioButtonModule,
    DialogModule,
    GrowlModule,
    TableModule
  ],
  providers: [
    HttpService,
    ProfileService
  ]
})
export class CreateProfileModule { }