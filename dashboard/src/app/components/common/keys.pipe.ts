import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: "Keys", pure: false })
export class Keys implements PipeTransform{
    transform(value: any, args: any[] = null): any{
        return Object.keys(value);
    }
}