import { Injectable, Output, EventEmitter } from '@angular/core';
import { Http } from '@angular/http';

@Injectable()
export class AppService {
    constructor(private http: Http){}

    @Output() onHeaderTitleChange = new EventEmitter<boolean>();

    changeHeaderTitle(){
        this.onHeaderTitleChange.emit();
    }

    logOut(){
        // return this.http.get("v1/portal/logout");
    }

}