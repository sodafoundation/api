import { DOCUMENT } from "@angular/common";

export class Mask {
    static createBackground = () =>{
        let div = document.createElement('div');
        div.className = 'mask-background';
        return div;
    }

    static createLoading= () => {
        let div = document.createElement('div');
        div.innerHTML = '<div class="M-loading"><div class="M-load M-loader"></div><div class="M-load M-loader1"></div><div class="M-load M-loader2"></div><div class="M-load M-loader3"></div><div class="M-load M-loader4"></div><div class="M-load M-loader5"></div></div><p class="M-loading-text">Loading...</p>';
        div.className = "M-loading-content";
        return div;
    }

    static background = Mask.createBackground();

    static loading = Mask.createLoading();

    static show() {
        document.body.appendChild(Mask.loading);
        document.body.appendChild(Mask.background);
    }

    static hide() {
        if( document.body.querySelector('body > .M-loading-content') ){
            document.body.removeChild(Mask.loading);
        }

        if( document.body.querySelector('body > .mask-background') ){
            document.body.removeChild(Mask.background);
        }
    }
}