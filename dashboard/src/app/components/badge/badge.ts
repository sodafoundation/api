import {NgModule,Directive, ElementRef,OnInit,AfterViewInit,OnDestroy,Input,Output,SimpleChange,EventEmitter,forwardRef,Renderer} from '@angular/core';
import {CommonModule} from '@angular/common';
import { DomHandler } from '../dom/domhandler';

@Directive({
    selector: '[pBadge]',
    providers:[DomHandler]
})
export class Badge implements OnInit,AfterViewInit{
    @Input('animate') animate:boolean = true;

    public _show:boolean = true;

    @Input('show') get show():boolean{
        return this._show;
    }

    set show(val:boolean){
        this._show = val;
        this.toggle();
        this.onBadgeChange.emit({'show':this.show,'count':this.count,'overflowCount':this.overflowCount});
    }

    public _count:number;

    get count():number{
        return this._count;
    }

    @Input('pBadge') set count(val:number){
        this._count = val>=0?val:0;
        this.onCountChange(val);
        this.onBadgeChange.emit({'show':this.show,'count':this.count,'overflowCount':this.overflowCount});
    }

    @Input('dot') dot:boolean = false;

    @Input() overflowCount:number = 99;

    @Input() showZero:boolean = false;

    @Input('status') status:string;//Enum{ 'success', 'processing, 'default', 'error', 'warning' }

    @Input('text') text:string;

    @Input('style') style:any;

    @Output() onBadgeChange:EventEmitter<any> = new EventEmitter();

    public container:any;

    constructor(public el: ElementRef,private domHnadler: DomHandler, public renderer: Renderer) {}

    ngOnInit(){
        if(this.status){
            this.dot = true;
        }
        this.createContainer();
    }

    ngAfterViewInit(){
        this.el.nativeElement.style.position = 'relative';
        this.el.nativeElement.style.overflow = 'visible';
        if(this.style){
            this.updateContainerStyle(this.style);
        }
    }

    createContainer(){
        if(this.count == 0 && !this.showZero) return;
        if(this.container) this.el.nativeElement.removeChild(this.container);
        this.container = document.createElement('span');
        this.container.className = 'ui-badge ui-widget';
        
        if(!this.show){
            this.container.classList.add('ui-badge-hide');
        }

        if(this.dot){
            this.container.classList.add('ui-badge-dot');
            if(this.status){
                this.updateStatus();
                return;
            }
        } else {
            this.container.title = '' + this.count || '';
        }

        this.updateContainerStyle(this.style);
        this.updateCountsHtml();
        this.el.nativeElement.appendChild(this.container);
    }

    updateContainerStyle(style:any){
        if(!this.container || typeof style !== 'object') return;
        for(let key in style){
            this.container.style[key] = style[key];
        }
    }

    updateCountsHtml(){
        if(!this.container || this.dot) return;
        if(this.animate && this.count <= this.overflowCount){
            if(this.domHnadler.find(this.container,'.ui-badge-counts').length > 0) return;

            let countsContainer = document.createElement('span');
            countsContainer.className = 'ui-badge-counts';
            let overflowArr = (this.overflowCount + '').split('');
            let countArr = (this.count + '').split('');
            for(let i = 0;i<overflowArr.length;i++){
                let countsItem = document.createElement('span');
                countsItem.className = 'ui-badge-counts-item';
                countsItem.style.transform = 'translateY(-' + 10 * parseInt(countArr[i]) + '%)';
                if(!countArr[i]) countsItem.style.display = 'none';
                for(let j = 0;j < 10;j++){
                    countsItem.innerHTML += '<p class="' + ((countArr[i] && countArr[i] == j + '')?'active':'') + '">' + j + '</p>';
                }
                countsContainer.appendChild(countsItem);
            }
            this.container.innerHTML = '';
            this.container.appendChild(countsContainer);
        } else {
            this.container.innerText = this.count > this.overflowCount?this.overflowCount + '+':this.count + '';
        }
    }

    updateStatus(){
        if(!this.container || !this.status) return;
        
        this.container.classList.add('ui-badge-status-content');
        this.container.innerHTML = '';
        this.container.innerHTML += '<span class="ui-badge-status ui-badge-status-' + this.status + '"></span>';
        
        if(this.text){
            this.container.innerHTML += '<span class="ui-badge-status-text">' + this.text + '</span>';
        }
        this.el.nativeElement.appendChild(this.container);
    }

    updateCountsPosition(){
        if(!this.animate || this.count > this.overflowCount || this.dot) return;
        let overflowArr = (this.overflowCount + '').split('');
        let countArr = (this.count + '').split('');
        let countsContainer = this.domHnadler.find(this.container,'.ui-badge-counts')[0];
        let counteItems = countsContainer.childNodes;
        for(let i = 0;i < overflowArr.length;i++){
            if(countArr[i]){
                counteItems[i].style.display = 'inline-block';
                counteItems[i].style.transform = 'translateY(-' + 10 * parseInt(countArr[i]) + '%)';
                for(let j = 0;j < 10;j++){
                    counteItems[i].childNodes[j].classList.remove('active');
                    if(countArr[i] == j + ''){
                        counteItems[i].childNodes[j].classList.add('active');
                    }
                }
            } else {
                counteItems[i].style.display = 'none';
                counteItems[i].style.transform = 'translateY(0%)';
            }
            
        }
    }

    onCountChange(count:number):void{
        if(this.dot || this.status) return;
        if(!this.container){
            this.createContainer();
            return;
        }
        this.container.title = '' + this.count || '';
        if(this.count < this.overflowCount && this.animate){
            this.updateCountsHtml();
            this.updateCountsPosition();
        } else {
            this.updateCountsHtml();
        }
        
    }

    toggle(){
        if(!this.container) return;
        if(this.show){
            this.container.classList.remove('ui-badge-hide');
        } else {
            this.container.classList.add('ui-badge-hide');
        }
    }

}

@NgModule({
    imports: [CommonModule],
    exports: [Badge],
    declarations: [Badge]
})
export class BadgeModule { }