import { NgModule, ModuleWithProviders, APP_INITIALIZER, Injector } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ExceptionService, MsgBoxService, I18NService, HttpService } from './api';
import { SharedConfig } from './shared.config';
import { I18N } from '../components/common/api';
import { XHRBackend, RequestOptions, Http } from '@angular/http';

export function httpFactory(backend: XHRBackend, options: RequestOptions, injector: Injector){
    options.headers.set("contentType", "application/json; charset=UTF-8");
    options.headers.set('Cache-control', 'no-cache');
    options.headers.set('cache-control', 'no-store');
    options.headers.set('expires', '0');
    options.headers.set('Pragma', 'no-cache');
    
    if( localStorage['x-subject-token'] ){
        options.headers.set('X-Auth-Token', localStorage['x-subject-token']);
    }

    return new HttpService(backend, options, injector);
}

@NgModule({
    exports:[FormsModule]
})

export class SharedModule {
    static forRoot(): ModuleWithProviders {
        return {
            ngModule: SharedModule,
            providers: [
                MsgBoxService,
                I18NService,
                ExceptionService,
                {
                    provide: Http,
                    useFactory: httpFactory,
                    deps: [XHRBackend, RequestOptions, Injector]
                },
                {
                    provide: APP_INITIALIZER,
                    useFactory: SharedConfig.config,
                    deps: [I18NService, Injector],
                    multi: true
                }
            ]
        };
    }
}