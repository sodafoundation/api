import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener,Input } from '@angular/core';
import { Http } from '@angular/http';

@Component({
    selector: 'img-item',
    templateUrl: './imgItem.component.html',
    styleUrls: [
        './imgItem.component.scss'
    ]
})
export class ImgItemComponent implements OnInit{
    @Input() item;

    
    constructor(
        private http: Http
    ){}
    
    ngOnInit() {
    }
}
