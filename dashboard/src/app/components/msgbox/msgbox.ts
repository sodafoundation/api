import { NgModule, Component, Input, Pipe } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { DialogModule } from '../dialog/dialog';

@Pipe({ name: 'safeHtml'})
export class Safe{
    constructor(private sanitizer: DomSanitizer){}

    transform (style): SafeHtml {
        return this.sanitizer.bypassSecurityTrustHtml(style);
    }
}

@Component({
    template:`
    <p-dialog [header]="config.header" [isMsgBox]="true" [modal]="true" [(visible)]="config.visible" [width]="config.width" [height]="config.height" (onOk)="config.ok()">
        <div class="msgbox">
            <div>
                <i [ngClass]="{'fa ':true, 'fa-info-circle': config.type=='info', 'error fa-times-circle': config.type=='error', 'success fa-info-circle': config.type=='success'}" [class]="config.icon" ></i>
            </div>
            <div>
                <h3 *ngIf="config.type=='success' || config.type=='error'" [ngClass]="config.type + '-title'">{{config.title}}</h3>
                <p [innerHTML]="config.content | safeHtml"></p>
            </div>
        </div>
    </p-dialog>`
})

export class MsgBox{
    @Input() config;
}

@NgModule({
    imports: [DialogModule, CommonModule],
    entryComponents: [MsgBox],
    exports: [MsgBox],
    declarations: [MsgBox, Safe]
})

export class MsgBoxModule{}