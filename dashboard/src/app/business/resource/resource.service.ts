import { Injectable } from '@angular/core';
import { I18NService, HttpService, ParamStorService } from '../../shared/api';
import { Observable } from 'rxjs';

@Injectable()
export class  AvailabilityZonesService{
  constructor(
    private http: HttpService,
    private paramStor: ParamStorService
  ) { }

  url = 'v1beta/{project_id}/availabilityZones';

  //get az
  getAZ(param?): Observable<any>{
    return this.http.get(this.url, param);
  }

  
}
