import { Injectable, Injector, ApplicationRef, ViewContainerRef, ComponentFactoryResolver } from '@angular/core';
import { MsgBox } from '../../components/msgbox/msgbox';
import { I18NService } from './I18N.service';

@Injectable()
export class MsgBoxService {
    private defaultConfig = {
        header: this.I18N.get("tip"),
        visible: true,
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

        let viewContainerRef = this.applicationRef.components[0].instance.viewContainerRef;
        let componentFactory = this.componentFactoryResolver.resolveComponentFactory(MsgBox);
        viewContainerRef.clear();
        let componentRef = viewContainerRef.createComponent(componentFactory);

        if( !config.ok ){
            config.ok = () => config.visible = false;
        }

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
            content: "",
            header: this.I18N.get("success")
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
            content: "",
            header: this.I18N.get("error")
        }

        return this.open(Object.assign({}, this.defaultConfig, _config, config));
    }

}