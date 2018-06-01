import {NgModule,Component,ElementRef,AfterViewInit,OnDestroy,Input,Output,Renderer2,Inject,forwardRef,ViewChild,AfterViewChecked,ContentChildren,EventEmitter,QueryList,ChangeDetectorRef} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DomHandler} from '../dom/domhandler';
import {MenuItem} from '../common/menuitem';
import {Location} from '@angular/common';
import {trigger,state,style,transition,animate} from '@angular/animations';
import {ButtonModule} from '../button/button';
import {Router} from '@angular/router';
import {RouterModule} from '@angular/router';

@Component({
    selector: 'p-dropMenuSub',
    template: `
        <ul [ngClass]="{'ui-widget-content ui-corner-all ui-submenu-list ui-shadow':!root}" class="ui-menu-list" (click)="listClick($event)">
            <ng-template ngFor let-child [ngForOf]="(root ? item : item.items)">
                <li *ngIf="child.separator" class="ui-menu-separator ui-widget-content">
                <li *ngIf="!child.separator" #item [ngClass]="{'ui-menuitem ui-corner-all':true,'ui-menuitem-active':item==activeItem}"
                    (mouseenter)="onItemMouseEnter($event,item,child)" (mouseleave)="onItemMouseLeave($event,item)">
                    <a *ngIf="!child.routerLink" [href]="child.url||'#'" [attr.target]="child.target" [attr.title]="child.title" [attr.id]="child.id" (click)="itemClick($event, child)"
                        [ngClass]="{'ui-menuitem-link ui-corner-all':true,'ui-state-disabled':child.disabled}" [ngStyle]="child.style" [class]="child.styleClass">
                        <span class="ui-submenu-icon fa fa-fw fa-caret-right" *ngIf="child.items"></span>
                        <span class="ui-menuitem-icon fa fa-fw" *ngIf="child.icon" [ngClass]="child.icon"></span>
                        <span class="ui-menuitem-text">{{child.label}}</span>
                    </a>
                    <a *ngIf="child.routerLink" [routerLink]="child.routerLink" [queryParams]="child.queryParams" [routerLinkActive]="'ui-state-active'" 
                        [routerLinkActiveOptions]="child.routerLinkActiveOptions||{exact:false}" [attr.target]="child.target" [attr.title]="child.title" [attr.id]="child.id"
                        (click)="itemClick($event, child)" [ngClass]="{'ui-menuitem-link ui-corner-all':true,'ui-state-disabled':child.disabled}" 
                        [ngStyle]="child.style" [class]="child.styleClass">
                        <span class="ui-submenu-icon fa fa-fw fa-caret-right" *ngIf="child.items"></span>
                        <span class="ui-menuitem-icon fa fa-fw" *ngIf="child.icon" [ngClass]="child.icon"></span>
                        <span class="ui-menuitem-text">{{child.label}}</span>
                    </a>
                    <p-dropMenuSub class="ui-submenu" [item]="child" *ngIf="child.items"></p-dropMenuSub>
                </li>
            </ng-template>
        </ul>
    `,
    providers: [DomHandler]
})
export class ContextMenuSub {

    @Input() item: MenuItem;
    
    @Input() root: boolean;
    
    constructor(public domHandler: DomHandler, @Inject(forwardRef(() => DropMenu)) public contextMenu: DropMenu) {}
        
    activeItem: any;

    containerLeft: any;

    hideTimeout: any;
                
    onItemMouseEnter(event, item, menuitem) {
        if(menuitem.disabled) {
            return;
        }

        if(this.hideTimeout) {
            clearTimeout(this.hideTimeout);
            this.hideTimeout = null;
        }
        
        this.activeItem = item;
        let nextElement =  item.children[0].nextElementSibling;
        if(nextElement) {
            let sublist = nextElement.children[0];
            sublist.style.zIndex = ++DomHandler.zindex;
            this.position(sublist, item);
        }
    }
    
    onItemMouseLeave(event, link) {
        this.hideTimeout = setTimeout(() => {
            this.activeItem = null;
        }, 1000);
    }
    
    itemClick(event, item: MenuItem) {
        if(item.disabled) {
            event.preventDefault();
            return;
        }
        
        if(!item.url) {
            event.preventDefault();
        }
        
        if(item.command) {            
            item.command({
                originalEvent: event,
                item: item
            });
        }
    }
    
    listClick(event) {
        this.activeItem = null;
    }
    
    position(sublist, item) {
        this.containerLeft = this.domHandler.getOffset(item.parentElement)
        let viewport = this.domHandler.getViewport();
        let sublistWidth = sublist.offsetParent ? sublist.offsetWidth: this.domHandler.getHiddenElementOuterWidth(sublist);
        let itemOuterWidth = this.domHandler.getOuterWidth(item.children[0]);

        sublist.style.top = '0px';

        if((parseInt(this.containerLeft.left) + itemOuterWidth + sublistWidth) > (viewport.width - this.calculateScrollbarWidth())) {
            sublist.style.left = -sublistWidth + 'px';
        }
        else {
            sublist.style.left = itemOuterWidth + 'px';
        }
    }

    calculateScrollbarWidth(): number {
        let scrollDiv = document.createElement("div");
        scrollDiv.className = "ui-scrollbar-measure";
        document.body.appendChild(scrollDiv);

        let scrollbarWidth = scrollDiv.offsetWidth - scrollDiv.clientWidth;
        document.body.removeChild(scrollDiv);
        
        return scrollbarWidth;
    }
}

@Component({
    selector: 'p-dropmenu',
    template: `
        <div #container [ngClass]="{'ui-dropmenu ui-buttonset ui-widget':true,'ui-state-disabled':disabled}" [ngStyle]="style">
            <button type="button" [label]="label" pButton [icon]="icon" iconPos="right" (click)="onDropdownButtonClick($event)" [disabled]="disabled"></button>
            <div #overlay [ngClass]="'ui-dropmenu-list ui-widget ui-widget-content ui-corner-all ui-shadow'" [class]="styleClass">
                <p-dropMenuSub [item]="model" root="root"></p-dropMenuSub>
            </div>
        </div>
    `,
    providers: [DomHandler]
})
export class DropMenu implements AfterViewInit,OnDestroy {

    @Input() model: MenuItem[];

    @Input() label: string;

    @Input() icon: string = "fa-caret-down";

    @Input() disabled: boolean;
    
    @Input() global: boolean;
    
    @Input() target: any;

    @Input() style: any;

    @Input() styleClass: string;
    
    @Input() appendTo: any = "body";

    @Input() autoZIndex: boolean = true;
    
    @Input() baseZIndex: number = 0;

    @Output() onDropdownClick: EventEmitter<any> = new EventEmitter();
    
    @ViewChild('container') containerViewChild: ElementRef;

    @ViewChild('overlay') overlayViewChild: ElementRef;
                    
    documentClickListener: any;

    public dropdownClick: boolean;

    public shown: boolean;
        
    constructor(public el: ElementRef, public domHandler: DomHandler, public renderer: Renderer2) {}

    ngAfterViewInit() {        
        if(this.appendTo) {
            if(this.appendTo === 'body')
                document.body.appendChild( this.overlayViewChild.nativeElement);
            else
                this.domHandler.appendChild(this.overlayViewChild.nativeElement, this.appendTo);
        }
    }
    
    onDropdownButtonClick(event: Event) {
        if(!this.shown) {
            this.dropdownClick = true;
        }else{
            this.dropdownClick = false;
        }
        this.onDropdownClick.emit(event);
        this.show(event);
    }
        
    show(event?) {
        this.alignPanel(); 
        this.moveOnTop();
        this.overlayViewChild.nativeElement.style.display = 'block';
        this.domHandler.fadeIn(this.overlayViewChild.nativeElement, 250);
        this.bindDocumentClickListener();
        
        if(event) {
            event.preventDefault();
        }
    }
    
    hide() {
        this.overlayViewChild.nativeElement.style.display = 'none';
        this.unbindDocumentClickListener();
    }

    alignPanel() {
        if(this.appendTo)
            this.domHandler.absolutePosition(this.overlayViewChild.nativeElement, this.containerViewChild.nativeElement);
        else
            this.domHandler.relativePosition(this.overlayViewChild.nativeElement, this.containerViewChild.nativeElement);
    }

    moveOnTop() {
        if(this.autoZIndex) {
            this.overlayViewChild.nativeElement.style.zIndex = String(this.baseZIndex + (++DomHandler.zindex));
        }
    }
    
    bindDocumentClickListener() {
        if(!this.documentClickListener) {
            this.documentClickListener = this.renderer.listen('document', 'click', (event) => {
                if(this.dropdownClick) {
                    this.dropdownClick = false;
                    this.shown = true;
                }
                else {
                    this.hide();
                    this.shown = false;
                    this.unbindDocumentClickListener();
                }
            });
        }
    }
    
    unbindDocumentClickListener() {
        if(this.documentClickListener) {
            this.documentClickListener();
            this.documentClickListener = null;
        }
    }
        
    ngOnDestroy() {
        this.unbindDocumentClickListener();
                
        if(this.appendTo) {
            this.el.nativeElement.appendChild(this.overlayViewChild.nativeElement);
        }
    }

}

@NgModule({
    imports: [CommonModule, ButtonModule, RouterModule],
    exports: [DropMenu, ButtonModule, RouterModule],
    declarations: [DropMenu,ContextMenuSub]
})
export class DropMenuModule { }