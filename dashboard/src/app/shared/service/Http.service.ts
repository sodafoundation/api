import { Http, XHRBackend, RequestOptions, Request, RequestOptionsArgs, Response, Headers, BaseRequestOptions } from '@angular/http';
import { Injectable, Injector } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Mask } from '../utils/Mask';
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
        // this.setToken(options);
        if( localStorage['x-subject-token'] ){
            !options && (options = {})
            !options.headers && (options['headers'] = new Headers());
            options.headers.set('X-Auth-Token', localStorage['x-subject-token']);
        }
        return this.intercept(super.get(url, options), options);
    }

    post(url: string, body: any, options?: BaseRequestOptionsArgs): Observable<Response>{
        return this.intercept(super.post(url, body, options), options);
    }

    put(url: string, body: any, options?: BaseRequestOptionsArgs): Observable<Response>{
        // this.setToken(options);
        if( localStorage['x-subject-token'] ){
            !options && (options = {})
            !options.headers && (options['headers'] = new Headers());
            options.headers.set('X-Auth-Token', localStorage['x-subject-token']);
        }
        return this.intercept(super.put(url, body, options), options);
    }

    delete(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        // this.setToken(options);
        if( localStorage['x-subject-token'] ){
            !options && (options = {})
            !options.headers && (options['headers'] = new Headers());
            options.headers.set('X-Auth-Token', localStorage['x-subject-token']);
        }
        return this.intercept(super.delete(url, options), options);
    }

    patch(url: string, body: any, options?: BaseRequestOptionsArgs): Observable<Response>{
        // this.setToken(options);
        if( localStorage['x-subject-token'] ){
            !options && (options = {})
            !options.headers && (options['headers'] = new Headers());
            options.headers.set('X-Auth-Token', localStorage['x-subject-token']);
        }
        return this.intercept(super.patch(url, body, options), options);
    }

    head(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        return this.intercept(super.head(url, options), options);
    }

    options(url: string, options?: BaseRequestOptionsArgs): Observable<Response>{
        return this.intercept(super.options(url, options), options);
    }

    intercept(observable: Observable<Response>, options = <BaseRequestOptionsArgs>{}): Observable<Response> {
        let exception = this.injector.get(ExceptionService);

        if(options.mask != false) {
            Mask.show();
        }

        return observable.timeout(options.timeout || this.TIMEOUT).do((res: Response) => {
            //success

            //fail: 公共异常处理
            if(exception.isException(res)){
                options.doException !== false && exception.doException(res);

                throw Observable.throw(res);
            }
        }, (err:any) => {
            //fail: 公共异常处理
            options.doException !== false && exception.doException(err);
        })
        .finally(() => {
            if( options.mask !==false ){
                Mask.hide();
            }
        })
    }
}