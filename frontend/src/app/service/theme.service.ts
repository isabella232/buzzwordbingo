import { Injectable } from '@angular/core';
import { BrowserModule, Title }  from '@angular/platform-browser'

export const unicornTheme = {
    'body-font':'"Google Sans", "Roboto", sans-serif',
    'header-font':'"Cinzel Decorative", cursive',
    'login-header-font':'"Roboto"',
    'header-before':'"🌈Unicorn "',
    'header-after':'"🦄"',
    'body-bg-color':'rgb(52, 47, 55)',
    'body-bg':'url("assets/wallpaper.jpg")',
    'bingo-header-bg-image':'url("assets/celebration.jpeg")',
    'interface-bg-color':'#260252',
    'interface-text-color':'#FFF',
    'link-text-color':'#fff25e',
    'link-text-bg-color':'#000',
    'link-hover-bg-color':'#F8F9FA',
    'link-hover-color':'#111',
    'link-active-bg-color':'#333',
    'link-active-color':'#7d7a5b',
    'bingo-header-text-shadow':'#260252',
    'bingo-header-color':'#e3c3f',
    'bingo-bg-color':'#ffedfb',
    'bingo-color':'#270147',
    'border-color':'#ccc',
    'login-bg-color':'#f5ceec',
    'login-highlight-bg-color':'#f222ad',
    'login-highlight-contrast-bg-color':'#82aff7',
    'login-alert-bg-color':'#ff0000',
    'login-readonly-bg-color':'#f5ceec',
    'login-readonly-color':'#666',
    'disabled-item-bg-color':'#474545',
    'disabled-item-color':'#777',
    'record-bg-color':'#f5e389',
    'record-color':'#c18a00',
    "input-bg-color":"#fff",
    "input-color":"#121212",
    "input-border-color":"#000",
    "bingo-declared-header-bg-color":"#260252",
    "bingo-declared-header-color":"#fff",
    "login-color":"#121212",
    "login-header-text-color":" #fff",
    'header-color':'#fff'
} 

export const lightTheme = {
  'body-font':'"Google Sans", "Roboto", sans-serif',
  'header-font':'"Google Sans", "Roboto", sans-serif',
  'login-header-font':'"Roboto"',
  'header-before':'""',
  'header-after':'""',
  'body-bg-color':'#ffffff',
  'body-bg':'none',
  'bingo-header-bg-image':'none',
  'interface-bg-color':'#eee',
  'interface-text-color':'#000',
  'link-text-color':'#03DAc6',
  'link-text-bg-color':'#fff',
  'link-hover-bg-color':'#03DAc6',
  'link-hover-color':'#121212',
  'link-active-bg-color':'#BBB',
  'link-active-color':'#121212',
  'bingo-header-text-shadow':'#026056',
  'bingo-bg-color':'#03DAc6',
  'bingo-color':'rgb(255, 255, 255, .80)', 
  'bingo-header-color':'#03DAc6',
  'border-color':'#ccc',
  'login-bg-color':'#fff',
  'login-color':'#000',
  'login-highlight-bg-color':'#6200EE',
  'login-highlight-contrast-bg-color':'#82aff7',
  'login-alert-bg-color':'#CF6679',
  'login-readonly-bg-color':'#ccc',
  'login-readonly-color':'#999',
  'disabled-item-bg-color':'#dedede',
  'disabled-item-color':'#777',
  'record-bg-color':'#03DAc6',
  'record-color':'rgb(255, 255, 255, .80)',  
  'input-bg-color':'#f9f9f9',
  'input-color':'#121212',
  'input-border-color':'#ccc',
  'bingo-declared-header-bg-color':'#03DAc6',
  'bingo-declared-header-color':' #fff',
  'login-header-text-color':'#fff',
  'header-color':'#121212'
} 


export const darkTheme = {
  'body-font':'"Google Sans", "Roboto", sans-serif',
  'header-font':'"Google Sans", "Roboto", sans-serif',
  'login-header-font':'"Roboto"',
  'header-before':'""',
  'header-after':'""',
  'body-bg-color':'#000000',
  'body-bg':'none',
  'bingo-header-bg-image':'none',
  'interface-bg-color':'#232323',
  'interface-text-color':'rgb(255, 255, 255, .87)',
  'link-text-color':'#03DAc6',
  'link-text-bg-color':'#232323',
  'link-hover-bg-color':'#03DAc6',
  'link-hover-color':'#232323',
  'link-active-bg-color':'#333',
  'link-active-color':'rgb(255, 255, 255, .60)',
  'bingo-header-text-shadow':'#026056',
  'bingo-bg-color':'#03DAc6',
  'bingo-color':'#270147',
  'bingo-header-color':'#03DAc6',
  'border-color':'#ccc',
  'login-bg-color':'#232323',
  'login-color':'rgb(255, 255, 255, .87)',
  'login-highlight-bg-color':'#BB86FC',
  'login-highlight-contrast-bg-color':'#82aff7',
  'login-alert-bg-color':'#CF6679',
  'login-readonly-bg-color':'#444',
  'login-readonly-color':'#bbb',
  'disabled-item-bg-color':'#474545',
  'disabled-item-color':'#777',
  'record-bg-color':'#BB86FC',
  'record-color':'rgb(255, 255, 255, .60)',  
  'input-bg-color':'#232323',
  'input-color':'rgba(255, 255, 255, .60)',
  'input-border-color':'#000',
  'bingo-declared-header-bg-color':'#03DAc6',
  'bingo-declared-header-color':' #232323',
  'header-color':'#fff'
} 








@Injectable({ providedIn: 'root' })
export class ThemeService {

  constructor(private titleService:Title) {
    this.toggleLight();

   }

  toggleDark() {
    this.setTheme(darkTheme);
    this.titleService.setTitle("Meeting Bingo");
  }

  toggleLight() {
    this.setTheme(lightTheme);
    this.titleService.setTitle("Meeting Bingo");
  }

  toggleUnicorn() {
    this.setTheme(unicornTheme);
    this.titleService.setTitle("🌈Unicorn Meeting Bingo🦄");
  }

  private setTheme(theme: {}) {
    console.log("Theme change called")
    Object.keys(theme).forEach(k =>
      document.documentElement.style.setProperty(`--${k}`, theme[k])
    );
  }
}
