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
    <p-dialog [header]="config.header" [isMsgBox]="true" [(visible)]="config.visible" [width]="config.width" [height]="config.height" [cancelBtnVisible]="config.cancelBtnVisible" (onOk)="config.ok()" (onCancel)="config.cancel()"
    [okBtnDisabled]="config.okBtnDisabled" [closeBtnDisabled]="config.closeBtnDisabled" [cancelBtnDisabled]="config.cancelBtnDisabled" [btnFocus]="config.btnFocus" [showCloseBtn]="config.showCloseBtn">
        <div class="msgbox">
            <div>
                <svg class="icon icon-size-48" [ngSwitch]="config.type">
                    <use xlink:href="#icon-popup-success-48" *ngSwitchCase="'success'"></use>
                    <use xlink:href="#icon-popup-error-48" *ngSwitchCase="'error'"></use>
                    <use xlink:href="#icon-popup-info-48" *ngSwitchCase="'info'"></use>
                    <use xlink:href="#icon-popup-question-48" *ngSwitchCase="'confirm'"></use>
                    <use xlink:href="#icon-popup-warning-48" *ngSwitchCase="'warn'"></use>
                </svg>
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