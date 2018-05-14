import { Injectable } from '@angular/core';

@Injectable()
export class ParamStorService {
    // Set local storage
    setParam(key, value){
        localStorage[key] = value;
    }
    
    // Get local storage
    getParam(key){
        return localStorage[key];
    }

    // Current login user
    CURRENT_USER(param?){
        if(param){
            this.setParam('current-user', param);
        } else {
            return this.getParam('current-user');
        }
    }

    // User auth token
    AUTH_TOKEN(param?){
        if(param){
            this.setParam('auth-token', param);
        } else {
            return this.getParam('auth-token');
        }
    }

    // Current login user
    CURRENT_TENANT(param?){
        if(param){
            this.setParam('current-tenant', param);
        } else {
            return this.getParam('current-tenant');
        }
    }


}