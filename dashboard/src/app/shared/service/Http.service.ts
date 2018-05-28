import { Http, XHRBackend, RequestOptions, Request, RequestOptionsArgs, Response, Headers, BaseRequestOptions } from '@angular/http';
import { Injectable, Injector } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Mask } from '../utils/mask';
import { BaseRequestOptionsArgs } from './api';
import { ExceptionService } from './Exception.service';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/finally';
import 'rxjs/add/operator/timeout';
import 'rxjs/add/observable/throw';

@Injectable()
export class HttpService extends Http {
    TIMEOUT = 18000;
    
    constructor(backend: XHRBackend, defaultOptions: RequestOptions, private injector: Injector) {
        super(backend, defaultOptions);
    }

    get(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.get(url, options), options);
    }

    post(url: string, body: any, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.post(url, body, options), options);
    }

    put(url: string, body: any, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.put(url, body, options), options);
    }

    delete(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.delete(url, options), options);
    }

    patch(url: string, body: any, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.patch(url, body, options), options);
    }

    head(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.head(url, options), options);
    }

    options(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        [url, options]= this.presetURL(url, options);
        return this.intercept(super.options(url, options), options);
    }

    presetURL(url, options){
        // Configure token option
        if( localStorage['auth-token'] ){
            !options && (options = {})
            !options.headers && (options['headers'] = new Headers());
            options.headers.set('X-Auth-Token', localStorage['auth-token']);
            
        }

        // Configure "project_id" for url
        if(url.includes('{project_id}')){
            let project_id = localStorage['current-tenant'].split("|")[1];
            url = url.replace('{project_id}',project_id);
            
        }

        return [url, options];
    }

    intercept(observable: Observable<Response>, options = <BaseRequestOptionsArgs>{}): Observable<Response> {
        let exception = this.injector.get(ExceptionService);

        if(options.mask != false) {
            Mask.show();
        }

        return observable.timeout(options.timeout || this.TIMEOUT).do((res: Response) => {
            //success

            //fail: Public exception handling
            if(exception.isException(res)){
                options.doException !== false && exception.doException(res);

                throw Observable.throw(res);
            }
        }, (err:any) => {
            //fail: Public exception handling
            options.doException !== false && exception.doException(err);
        })
        .finally(() => {
            if( options.mask !==false ){
                Mask.hide();
            }
        })
    }
}