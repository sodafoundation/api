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

  //修改 volume
  modifyVolume(id,param) {
    let modifyUrl = this.url + '/' + id
    return this.http.put(modifyUrl, param);
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

  //查询指定Id volume
  getVolumeById(id): Observable<any> {
    let url = this.url + '/' + id;
    return this.http.get(url);
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

@Injectable()
export class SnapshotService {
  url = 'v1beta/ef305038-cd12-4f3b-90bd-0612f83e14ee/block/snapshots'
  constructor(private http: HttpService) { }
  //创建 snapshot
  createSnapshot(param) {
    return this.http.post(this.url, param);
  }

  //删除 snapshot
  deleteSnapshot(id){
    let url = this.url + "/" + id;
    return this.http.delete(url);
  }

  //查询 snapshot
  getSnapshots(filter?){
    let url = this.url;
    if(filter){
      url = this.url + "?" + filter.key + "=" + filter.value;
    }
    console.log(url);
    return this.http.get(url);
  }

  //修改 snapshot
  modifySnapshot(id,param){
    let url = this.url + "/" + id;
    return this.http.put(url,param);
  }
}
