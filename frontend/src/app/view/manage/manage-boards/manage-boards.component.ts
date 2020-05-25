import { Component, OnInit, Input } from '@angular/core';
import { Observable, of } from 'rxjs';
import { DataService, Phrase} from '../../../service/data.service'
import {GameService, Board, Message, Record} from '../../../service/game.service'

@Component({
  selector: 'app-manage-boards',
  templateUrl: './manage-boards.component.html',
  styleUrls: ['./manage-boards.component.scss']
})
export class ManageBoardsComponent implements OnInit {
  @Input() id:string;
  public boards: Observable<any[]>;
  constructor(public data:DataService, public game:GameService) { }

  ngOnInit(): void {
    this.boards = this.data.getBoards(this.id)
    this.boards.subscribe();
  }

  reset(bid:string, gid:string){
    console.log("Reset called")
    this.game.resetboard(bid, gid);
  }


}