import { writable } from 'svelte/store';

import unix from 'moment';

export const apiURL = "http://localhost:3000/"

export function isLoggedIn(){
    if (localStorage.getItem("token") == null){
        return false
    }
    return true
}

export function toDateString(timestamp: number){
    return unix(timestamp*1000).format("YYYY-MM-DD");
}