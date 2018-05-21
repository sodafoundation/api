import { EventEmitter, Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';

// export{ DomHandler } from '../dom/domhandler';

import{ RequestOptionsArgs } from '@angular/http';

export interface SortMeta {

}

export interface Confirmation {
    message: string;
    icon?: string;
    header?: string;
    accept?: Function;
    reject?: Function;
    acceptVisible?: boolean;
    rejectVisible?: boolean;
    acceptEvent?: EventEmitter<any>;
    rejectEvent?: EventEmitter<any>;
}

export interface BaseRequestOptionsArgs extends RequestOptionsArgs {
    mask?: boolean;
    timeout?: number;
    doException?: boolean;
}

@Injectable()
export class ConfirmationService{
    private requireConfirmationSource = new Subject<Confirmation>();
    private acceptConfirmationSource = new Subject<Confirmation>();

    requireConfirmation$ = this.requireConfirmationSource.asObservable();
    accept = this.acceptConfirmationSource.asObservable();

    confirm(confirmation: Confirmation) {
        this.requireConfirmationSource.next(confirmation);
        return this;
    }

    onAccept() {
        this.acceptConfirmationSource.next();
    }
}