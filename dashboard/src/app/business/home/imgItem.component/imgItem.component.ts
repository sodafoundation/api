import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener,Input } from '@angular/core';
import { Http } from '@angular/http';

@Component({
    selector: 'img-item',
    templateUrl: './imgItem.component.html',
    styleUrls: []
})
export class ImgItemComponent implements OnInit{
    countNum:number = 0;
    label:string = '';
    @Input() item;

    
    constructor(
        private http: Http
    ){}
    
    ngOnInit() {
        this.countNum = this.item.countNum;
        this.label = this.item.label;
    }
}
