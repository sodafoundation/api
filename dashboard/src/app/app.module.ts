import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CommonModule } from '@angular/common';
import { NgModule, APP_INITIALIZER } from '@angular/core';
import { HttpModule } from "@angular/http";
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SharedModule } from './shared/shared.module';
import { DropMenuModule, SelectButtonModule, ButtonModule, InputTextModule } from './components/common/api';
// import { AppService } from './app.service';
import { LocationStrategy, HashLocationStrategy } from '@angular/common';

import { MessagesModule } from './components/messages/messages';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    CommonModule,
    AppRoutingModule,
    MessagesModule,
    HttpModule,
    BrowserAnimationsModule,
    SharedModule.forRoot(),
    DropMenuModule,
    SelectButtonModule,
    ButtonModule,
    InputTextModule
  ],
  providers: [
    // AppService,
      { provide: LocationStrategy, useClass: HashLocationStrategy }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }