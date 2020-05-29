import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { share } from 'rxjs/operators';
import {GameService, Game} from '../service/game.service'
import {Router} from '@angular/router';
import { LocalstorageService } from './localstorage.service';

export class Player{
  name:string
  email:string
  admin:boolean
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private player:Player= new Player;
  private identity:any= new Player;
  private isAuthed:boolean=false;
  private isGameAdministrator:boolean=false;
  private isAdministrator:boolean=false;
  private hostUrl: string = environment.host_url;

  constructor(private http: HttpClient, 
              public game:GameService, 
              private localStorageService:LocalstorageService, 
              private router: Router) { 
    let player = localStorageService.getPlayer();
    if (player != null){
      this.setPlayer(player.name, player.email)
    }
  }
  

  setPlayer(name:string, email:string, admin:boolean=false){
    this.player.name = name;
    this.player.email = email;
    this.player.admin = admin;
    this.isAuthed = true;
    this.localStorageService.setPlayer(this.player);
    this.game.isAdmin().pipe(share()).subscribe(val=>{this.isAdministrator = val})

  } 

  getPlayer():Player{
    return this.player;
  }

  getIdentifiedEmail():string{
    return this.identity.email;
  }

  isAuth():boolean{
    return this.isAuthed;
  }

  isAdmin():boolean{
    return this.isAdministrator;
  }

  checkGameAdmin(gid:string){
    this.game.isGameAdmin(gid).pipe(share()).subscribe(val=>{this.isGameAdministrator = val})
  }

  isGameAdmin():boolean{
    return this.isGameAdministrator;
  }

  identifyPlayer () {
    return this.http.get(this.hostUrl + "/api/player/identify");
  }

  logout (reason:string="logged out") {
    console.log("logged out, reason:", reason )
    this.localStorageService.clearGameData();
    this.router.navigateByUrl('/login');
    return 
  }

}
