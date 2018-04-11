// import { Router } from '@angular/router';
import { Component, OnInit, ViewContainerRef, ViewChild, Directive, ElementRef, HostBinding, HostListener } from '@angular/core';
import { I18NService } from 'app/shared/api';
// import { AppService } from 'app/app.service';
import { I18nPluralPipe } from '@angular/common';
import { MenuItem } from './components/common/api';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: []
})
export class AppComponent implements OnInit{
    chromeBrowser: boolean = false;
    linkUrl = "";
    dropMenuItems: MenuItem[];
    menuItems = [
        {
            "title": "Home",
            "description": "Update 5 minutes ago",
            "routerLink": "/home"
        },
        {
            "title": "Block Service",
            "description": "23 volumes",
            "routerLink": "/block"
        },
        {
            "title": "profile",
            "description": "7 profiles have been created",
            "routerLink": "/profile"
        },
        {
            "title": "Resource",
            "description": "5 storages, 2 availability zone",
            "routerLink": "/resource"
        },
        {
            "title": "Identity",
            "description": "Managing tenants and users",
            "routerLink": "/identity"
        }
    ];
    activeItem: any;

    private msgs: any = [{ severity: 'warn', summary: 'Warn Message', detail: 'There are unsaved changes'}];

    constructor(
        private el: ElementRef
        // private I18N: I18NService,
        // private viewContainerRef: ViewContainerRef,
        // private appService: AppService,
        // private router: Router
    ){}
    
    ngOnInit() {
        this.activeItem = this.menuItems[0];

        this.dropMenuItems = [
            { 
                label: "Switch Tenant", 
                items:[
                    {
                        label: "TenantA", command:()=>{}
                    },
                    {
                        label: "TenantB", command:()=>{}
                    }
                ]
            },
            { label: "Logout", command:()=>{} }
        ];
    }
    
    menuItemClick(event, item) {
        this.activeItem = item;
    }

    supportCurrentBrowser(){
        let ie,
            firefox,
            safari,
            chrome,
            cIE = 11,
            cFirefox = 40,
            cChrome = 40;
        let ua = navigator.userAgent.toLowerCase();
        let isLinux = (ua.indexOf('linux') >= 0);

        if(this.isIE()) {
            if(ua.indexOf('msie') >= 0) {
                ie = this.getSys(ua.match(/msie ([\d]+)/));
            } else {
                ie = this.getSys(ua.match(/trident.*rv:([\d]+)/));
            }
        }else if(navigator.userAgent.indexOf("Firefox") > 0){
            firefox = this.getSys(ua.match(/firefox\/([\d]+)/));
        }else if(ua.indexOf("safari") != -1 && !(ua.indexOf("chrome") != -1)) {
            safari = this.getSys(ua.match(/version\/([\d]+)/));
        }else if(ua.indexOf("chrome") != -1) {
            chrome = this.getSys(ua.match(/chrome\/([\d]+)/));
        }

        if ((firefox) / 1 < cFirefox || (chrome) / 1 < cChrome || (ie) / 1 < cIE) {
            return true;
        }

        return false;
    }

    isIE() {
        return navigator.userAgent.toLowerCase().indexOf('trident') >= 0;
    }

    getSys (browserVersionArr) {
        if( !browserVersionArr) {
            return 0;
        } else if( browserVersionArr.length < 2) {
            return 0;
        } else {
            return browserVersionArr[1];
        }
    }
}
