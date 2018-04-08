import { Injectable, Injector, ApplicationRef, ViewContainerRef, ComponentFactoryResolver } from '@angular/core';
import { MsgBox } from '../../components/MsgBox/MsgBox';
import { I18NService } from './I18N.service';

@Injectable()
export class MsgBoxService {
    private defaultConfig = {
        header: this.I18N.get("tip"),
        visible: true,
        concelBtnVisible: false,
        closeBtnDisabled: false,
        okBtnDisabled: false,
        width: 418,
        height: "auto"
    }

    constructor(
        private componentFactoryResolver: ComponentFactoryResolver,
        private applicationRef: ApplicationRef,
        private I18N: I18NService,
        private injector: Injector
    ){}

    open(config){
        this.applicationRef = this.applicationRef || this.injector.get(ApplicationRef);
        
        if( !this.applicationRef || !this.applicationRef.components[0]){
            alert("System error");
            return;
        }

        let ViewContainerRef = this.applicationRef.components[0].instance.ViewContainerRef;
        let componentFactory = this.componentFactoryResolver.resolveComponentFactory(MsgBox);
        ViewContainerRef.clear();
        let componentRef = ViewContainerRef.createComponent(componentFactory);

        if( !config.ok ){
            config.ok = () => config.visible = false;
        }

        if( !config.cancel ){
            config.cancel = () => config.visible = false;
        }

        config.close = () => config.visible = false;

        (<MsgBox>componentRef.instance).config = config;

        return config;
    }

    info(config = {}){
        if( typeof config == "string"){
            config = {
                content: config
            }
        }

        let _config = {
            type: "info",
            content: "",
            header: this.I18N.get("tip"),
            btnFocus: "okBtn"
        }

        return this.open(Object.assign({}, this.defaultConfig, _config, config));
    }

    success(config = {}){
        if( typeof config == "string"){
            config = {
                content: config
            }
        }

        let _config = {
            type: "success",
            title: this.I18N.get("success"),
            content: "",
            header: this.I18N.get("tip")
        }

        return this.open(Object.assign({}, this.defaultConfig, _config, config));
    }

    error(config = {}){
        if( typeof config == "string"){
            config = {
                content: config
            }
        }

        let _config = {
            type: "error",
            title: this.I18N.get("fail"),
            content: "",
            header: this.I18N.get("tip"),
            btnFocus: "none"
        }

        return this.open(Object.assign({}, this.defaultConfig, _config, config));
    }

    confirm(config = {}){
        if( typeof config == "string"){
            config = {
                content: config
            }
        }

        let _config = {
            type: "confirm",
            title: this.I18N.get("confirm"),
            content: "",
            cancelBtnVisible: true,
            btnFocus: "none"
        }

        return this.open(Object.assign({}, this.defaultConfig, _config, config));
    }

    warn(config = {}){
        if( typeof config == "string"){
            config = {
                content: config
            }
        }

        let _config = {
            type: "warn",
            title: this.I18N.get("warn"),
            content: "",
            cancelBtnVisible: true,
            btnFocus: "none"
        }

        return this.open(Object.assign({}, this.defaultConfig, _config, config));
    }
}