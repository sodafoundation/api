import { Injectable } from '@angular/core';
import { HttpService } from './../../shared/service/Http.service';
import { Observable } from 'rxjs';

@Injectable()
export class ProfileService {
    url = 'v1beta/{project_id}/profiles'
    constructor(private http: HttpService) { }
    //创建profile
    createProfile(param) {
        return this.http.post(this.url, param);
    }

    //删除profile
    deleteProfile(id): Observable<any> {
        let deleteUrl = this.url + '/' + id
        return this.http.delete(deleteUrl);
    }

    //查询profiles
    getProfiles(): Observable<any> {
        return this.http.get(this.url);
    }

    //查询profiles
    getProfileById(id) {
        let getUrl = this.url + '/' + id
        return this.http.get(getUrl);
    }

    //修改profile
    modifyProfile(id, param) {
        let modifyUrl = this.url + '/' + id
        this.http.put(modifyUrl, param).subscribe((res) => {
            console.log(res.json().data);
        });
    }
}

@Injectable()
export class PoolService{
    url = 'v1beta/{project_id}/pools';
    constructor(private http: HttpService) { }
    //查询profiles
    getPools(): Observable<any> {
        return this.http.get(this.url);
    }
}
