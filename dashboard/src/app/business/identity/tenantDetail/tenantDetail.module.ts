import { NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TenantDetailComponent } from './tenantDetail.component';
import { ButtonModule, DataTableModule, DropMenuModule, DialogModule, ConfirmDialogModule, BadgeModule, InputTextModule } from '../../../components/common/api';

@NgModule({
  declarations: [ TenantDetailComponent ],
  imports: [ CommonModule, ButtonModule, DataTableModule, DropMenuModule, DialogModule, ConfirmDialogModule, BadgeModule, InputTextModule ],
  exports: [ TenantDetailComponent ],
  providers: []
})
export class TenantDetailModule { }