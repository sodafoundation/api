import { Injectable } from '@angular/core';
import { HttpService } from './../../shared/service/Http.service';
import { Observable } from 'rxjs';

@Injectable()
export class VolumeService {
  url = 'v1beta/ef305038-cd12-4f3b-90bd-0612f83e14ee/block/volumes'
  constructor(private http: HttpService) { }
  //创建 volume
  createVolume(param) {
    return this.http.post(this.url, param);
  }

  //删除 volume
  deleteVolume(id): Observable<any> {
    let deleteUrl = this.url + '/' + id
    return this.http.delete(deleteUrl);
  }

  //查询 volumes
  getVolumes(): Observable<any> {
    return this.http.get(this.url);
  }

  //创建 volumesGroup
  createVolumesGroup(param) {
    return this.http.post(this.url, param);
  }

  //删除 volumesGroup
  deleteVolumesGroup(id): Observable<any> {
    let deleteUrl = this.url + '/' + id
    return this.http.delete(deleteUrl);
  }

  //查询 volumesGroups
  getVolumesGroups(): Observable<any> {
    return this.http.get(this.url);
  }
}
