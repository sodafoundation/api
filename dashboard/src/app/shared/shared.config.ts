import{ Consts } from 'app/shared/api';
import{ Http } from '@angular/http';
import{ I18N } from 'app/components/common/api';

var _ = require('underscore');
export class SharedConfig{

    //配置入口
    static config( I18NService, injector) {
        let httpService = injector.get(Http);
        return () => new Promise((resolve, reject) => {
            Promise.all([
                SharedConfig.I18NConfig(I18NService, httpService),
                // SharedConfig.AutoDeploy(httpService),
            ]).then(() => resolve()).catch(reason => console.log(reason));
        })
    }

    //I18N服务配置
    static I18NConfig(I18NService, httpService) {
        let p1 = new Promise((resolve, reject) => {
            //应用I18N配置
            httpService.get("src/app/i18n/"+ I18NService.language + "/keyID.json").subscribe((r) => {
                Object.assign(I18NService.keyID, r.json());
                resolve();
            })
        })

        let p2 = new Promise((resolve, reject) => {
            //异常I18N配置
            httpService.get("src/app/i18n/"+ I18NService.language + "/exception.json").subscribe((r) => {
                Object.assign(I18NService.keyID, r.json());
                resolve();
            })
        })

        let p3 = new Promise((resolve, reject) => {
            //控件I18N配置
            httpService.get("src/app/i18n/"+ I18NService.language + "/widgetsKeyID.json").subscribe((r) => {
                Object.assign(I18NService.keyID, r.json());
                I18N.language = I18NService.language;
                I18N.keyID = I18NService.keyID;
                resolve();
            })
        })
        return Promise.all([p1, p2, p3]);
    }

    static AutoDeploy(httpService) {
        return new Promise((resolve, reject) => {
            //读取环境配置url
            httpService.get("src/prodconfig/deployConfig.json").subscribe((r) => {
                let responseData = r.json;

                Promise.all([
                    // SharedConfig.initBaseUser(httpService),
                    // SharedConfig.initUserRoles(httpService)
                ]).then(() => resolve()).catch(reason => console.log(reason));
            })
        })
    }

    //获取用户基本权限
    static initBaseUser(httpService){
        return new Promise((resolve, reject) => {
            let request: any = { params: {}};
            request.params.type = "base";
            httpService.get("/v1/...", request).subscribe(response => {
                let responseData = response.json().data;
                // Consts.CURRENT_USERNAME = responseData.userName;
                // Consts.BASE_TENANT = responseData.tenant;

                resolve();
            })
        })
    }

    //获取用户角色信息
    static initUserRoles(httpService){
        return new Promise((resolve, reject) => {
            let request: any = { params: {}};
            httpService.get("/v1/...", request).subscribe(response => {
                let roles = response.json().data;
                // _.each(roles, function(item, i){
                //     roles[i] = item.toLowerCase();
                // })
                // Consts.USER_ROLES = roles;

                resolve();
            })
        })
    }
}