import { Injectable } from '@angular/core';
import { templateJitUrl } from '@angular/compiler';

@Injectable()

export class I18NService {
    language = 'en';
    keyID = {};

    constructor(){
        this.language = this.getCookieLanguage();
    }

    get(key, params = []){
        let str = this.keyID[key] || key;
        params.forEach((param, index) => {
            str = str.replace('{' + index + '}', param);
        });
        return str;
    }

    //zh: Chinese en: English
    getCookieLanguage() {
        this.urlLanguage();
        let cookieStr = document.cookie.split("; ");
        let arrLength = cookieStr.length;
        for (let i = 0; i <arrLength; i++){
            let tempArr = cookieStr[i].split("=");
            if ("language" === tempArr[0]){
                return tempArr[1];
            }
        }

        let language = window.navigator.language ? window.navigator.language.split("-")[0] : "en";
        if(language == "zh" || language == "en"){
            return language;
        } else{
            return "en";
        }
    }

    urlLanguage() {
        let url = window.location.href;
        let hasUrlLan = url.split("/#/")[0].split("?locale=");
        if (hasUrlLan.length >0){
            let tempLanguage = hasUrlLan[1];
            if(tempLanguage == "zh" || tempLanguage == "en"){
                document.cookie = "language=" + tempLanguage + ";domain=.huawei.com;path=/";
            }
        }
    }
}