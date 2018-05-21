import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TenantListComponent } from './tenantList.component';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { FormModule, CheckboxModule, ConfirmDialogModule, ButtonModule, MultiSelectModule, DataTableModule, DropMenuModule, DialogModule, InputTextModule, InputTextareaModule, DropdownModule } from '../../components/common/api';
import { TenantDetailModule } from './tenantDetail/tenantDetail.module';

@NgModule({
  declarations: [ TenantListComponent ],
  imports: [ 
    CommonModule, 
    ButtonModule, 
    DataTableModule, 
    DropMenuModule, 
    DialogModule,
    InputTextModule,
    InputTextareaModule,
    ReactiveFormsModule,
    FormsModule,
    FormModule,
    ConfirmDialogModule,
    MultiSelectModule,
    CheckboxModule,
    TenantDetailModule 
  ],
  exports: [ TenantListComponent ],
  providers: []
})
export class TenantListModule { }