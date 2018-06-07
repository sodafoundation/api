import { CommonModule } from '@angular/common';
import { NgModule, APP_INITIALIZER } from '@angular/core';
import { ProfileComponent } from './profile.component';
import { RouterModule } from '@angular/router';

import { ProfileCardComponent } from './profileCard/profile-card.component';
import { ButtonModule,CardModule,ChartModule,MessageModule,OverlayPanelModule,DialogModule ,ConfirmationService,ConfirmDialogModule} from '../../components/common/api';
import { ProfileService } from './profile.service';
import { HttpService } from '../../shared/api';
import { SuspensionFrameComponent } from './profileCard/suspension-frame/suspension-frame.component';

let routers = [{
  path: '',
  component: ProfileComponent
}]

@NgModule({
  declarations: [
    ProfileComponent,
    ProfileCardComponent,
    SuspensionFrameComponent
  ],
  imports: [
    RouterModule.forChild(routers),
    ButtonModule,
    CommonModule,
    CardModule,
    ChartModule,
    MessageModule,
    OverlayPanelModule,
    DialogModule,
    ConfirmDialogModule
  ],
  providers: [
    HttpService,
    ProfileService,
    ConfirmationService
  ]
})
export class ProfileModule { }