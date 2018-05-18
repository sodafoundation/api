import { Injectable } from '@angular/core';
import { I18NService } from './I18N.service';
import { MsgBoxService } from './MsgBox.service';


@Injectable()
export class ExceptionService {
    constructor(private I18N: I18NService, private msgBox: MsgBoxService) {}

    //异常判断
    isException(res){
        return false;
        // let isRest = res.url.indexOf('v1') || res.url.indexOf('v3');
        // if (isRest > 0){
        //     if(res.text() === "sessionout"){
        //         return true;
        //     }
        //     let data = res.json();
        //     if(data && data.code != 0) {
        //         return true;
        //     }
        //     else{
        //         return false;
        //     }
        // }
    }

    doException(res){
        // if(res.text() === "sessionout" ) {
        //     let url = window.location.href;
        //     let isresource = url.indexOf('siteDistribution');
        //     if(isresource > 0 ){
        //         parent.postMessage({reload: "reload" }, '*');
        //     }
        //     else{
        //         window.location.reload();
        //     }
        //     return;
        // }

        // let data = res.json();
        // let code = data.code;

        // if(data && data.code) {
        //     code = data.code;
        // }
        // this.msgBox.error(this.I18N.get(code) + "");
    }
}