export class I18N {
    static language = "en";
    static keyID = {};

    static get(key, params = []) {
        let str = this.keyID[key] || key;
        params.forEach((param, index) => {
            str = str.replace("{" + index + "}", param);
        });
        return str;
    }
}