import { NgModule, Component, Directive, Input, Output, ElementRef, Injector, OnInit, AfterViewInit, AfterContentChecked } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormGroup, AbstractControl } from '@angular/forms';
import { SharedModule } from '../../shared/shared.module';

@Directive({
    selector: 'form'
})

export class Form{
    @Input() grid: {label: string, content: string };

    @Input() errorMessage: {[key: string]: {[key: string]: any}};

    @Input() formGroup: FormGroup;

    constructor(el: ElementRef) {
        el.nativeElement.classList.add('form');
    }
}

@Component({
    selector: 'form-item',
    template: `
    <div class='form-item ui-g' *ngIf="!hide">
        <div [ngClass]='{"required": required, "form-label": true}' [class]='labelStyleClass'>
            <label>{{label}}</label>
        </div>
        <div [ngClass]='{"form-content": true}' [class]='contentStyleClass'>
            <ng-content></ng-content>
            <div class="ui-message ui-message-error ui-corner-all">
                <svg class="ui-message-error-icon icon">
                    <use xlink:href="#icon-status-failed"></use>
                </svg>
                <span class="ui-message-error-text" *ngFor="let errKey of [(formCtrl.errors | Keys)[0]]">
                    {{ (errorMessage && errorMessage[errKey]) || formCtrl.errors[errKey] }}
                </span>
            </div>
        </div>
    </div>
    `
})

export class FormItem implements OnInit, AfterViewInit, AfterContentChecked {

    @Input() required: boolean;

    @Input() label: string;

    @Input() hide: string;

    errorMessage: {[Key:string]:{[key:string]:any}};

    labelStyleClass: string;

    contentStyleClass: string;

    formInstance: Form;

    formCtrl: AbstractControl;

    formctrls: { name:string, formCtrl: AbstractControl, errorMessage: {[Key:string]:{[key:string]:any}}}[] = [];

    constructor( private el: ElementRef, private injector: Injector){}

    ngOnInit(): void{
        //栅格样式
        this.formInstance = this.injector.get(Form);
    
        if( !this.formInstance ){
            throw "require Form";
        }

        if( this.formInstance.grid ){
            this.labelStyleClass = this.formInstance.grid.label;
            this.contentStyleClass = this.formInstance.grid.content;
        }
    }

    ngAfterViewInit(): void{
        let ctrlElems =this.el.nativeElement.querySelectorAll(".form-content [formControlName]");
        if( this.formInstance.formGroup && ctrlElems.length > 0){
            ctrlElems.forEach(elem => {
                let name = elem.getattribute("formControlName");
                let formCtrl = this.formInstance.formGroup.get(name);

                if( formCtrl ){
                    this.formctrls.push({
                        name: name,
                        formCtrl: formCtrl,
                        errorMessage: this.formInstance.errorMessage && this.formInstance.errorMessage[name]
                    });
                }
            })
        }
    }

    ngAfterContentChecked(){
        this.updateFormCtrl();
    }

    //显示第一个错误
    updateFormCtrl(): void {
        let firstErrorCtrl = this.formctrls.filter( item => {
            let formCtrl = item.formCtrl;
            return formCtrl.invalid && (formCtrl.dirty || formCtrl.touched);
        })[0];

        if( firstErrorCtrl ){
            this.formCtrl = firstErrorCtrl.formCtrl;
            this.errorMessage = firstErrorCtrl.errorMessage;
        }
        else{
            this.formCtrl = null;
            this.errorMessage = null;
        }
    }

}

@NgModule({
    imports: [CommonModule, FormsModule, SharedModule],
    exports: [Form, FormItem],
    declarations: [Form, FormItem]
})
export class FormModule{}

