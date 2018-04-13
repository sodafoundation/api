import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { Http } from '@angular/http';

@Component({
    templateUrl: './home.component.html',
    styleUrls: []
})
export class HomeComponent implements OnInit{

    constructor(
        private http: Http
        // private I18N: I18NService,
        // private router: Router
    ){}
    
    ngOnInit() {
        let arr = [3,10,11,17,21,23,4];
        // let arr = [1,4,3,2];
        // this.arraySortUpdate(arr);

      
    }
    
    showData() {
        this.http.get("/v1beta/ef305038-cd12-4f3b-90bd-0612f83e14ee/profiles").subscribe((res)=>{
            console.log(res.json().data);
        });
    }
}
