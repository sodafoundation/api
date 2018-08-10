import { Component, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { CreateProfileComponent } from './createProfile.component';
import { RouterModule } from '@angular/router';
import { InputTextModule,InputTextareaModule, CheckboxModule, FormModule, ButtonModule, DropdownModule, RadioButtonModule, DialogModule, Message, GrowlModule, SelectButtonModule } from '../../../components/common/api';
import { HttpService } from './../../../shared/api';
import { ProfileService,PoolService } from './../profile.service';
import { PoolModule } from './../storage-pools-table/storage-pools-table.module';

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
    PoolModule,
    SelectButtonModule,
    FormModule,
    InputTextareaModule
  ],
  providers: [
    HttpService,
    ProfileService,
    PoolService
  ]
})
export class CreateProfileModule { }
