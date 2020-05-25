import { Component, OnInit, Input, Output, EventEmitter, ViewChild } from '@angular/core';
import { Observable, of } from 'rxjs';
import { Phrase} from '../../../service/data.service'
import {GameService, Board} from '../../../service/game.service'


@Component({
  selector: 'app-item',
  templateUrl: './item.component.html',
  styleUrls: ['./item.component.scss']
})

export class ItemComponent implements OnInit {
  @Input() phrase: Phrase;
  @Input() boardid: string;
  @Input() board: Board;
  @Input() currentState:any;
  @Input() position: number;
  @Input() bingo:boolean = false;
  @Output() phraseEmitter = new EventEmitter<Phrase>();
  @Output() readyEmitter = new EventEmitter<ItemComponent>();
  disabled:boolean=false;
  
  constructor(private game:GameService) { }

  ngOnInit(): void {
    this.readyEmitter.emit(this);
  }

  ngAfterViewChecked(): void {
    if (this.phrase.selected){
      this.setDisplayChecked();
    }
  }

  ngAfterViewInit():void{
    
  } 

  ngOnChange(){
    console.log("Bingo:", this.bingo);
  }

  select(){
    if (this.bingo){
      this.disabled = true;
      this.disable();
    }

    if (this.phrase.text == "FREE"){
      return;
    }

    if (this.disabled){
      return;
    }
    
    this.selectDisplay();
    this.phraseEmitter.emit(this.phrase);
    this.game.record(this.phrase.id, this.board.game, this.board.id);
  }


  setDisplayChecked(){
    let item:HTMLElement = document.querySelector("#id_"+ this.phrase.id);
    let selectedPhraseCount = Object.keys(this.currentState).length;
    // item.classList.add("selected");
    switch(selectedPhraseCount) {
    case 0:
      item.style.backgroundColor = '#e3c3ff';
      item.style.color = '#3f3d40';
      break;
    case 1:
      item.style.backgroundColor = '#facdd5';
      item.style.color = '#3f3d40';
      break;
    case 2:
      item.style.backgroundColor = '#fffcbc';
      item.style.color = '#3f3d40';
      break;
    case 3:
      item.style.backgroundColor = '#d9ffe3';
      item.style.color = '#3f3d40';
      break;
    case 4:
      item.style.backgroundColor = '#ccedfd';
      item.style.color = '#3f3d40';
      break;
    case 5:
      item.style.backgroundColor = '#e3c3ff';
      item.style.color = '#3f3d40';
      break;


    default:
      item.style.backgroundColor = '#e3c3ff';
      item.style.color = '#3f3d40';
      break;

    }
  }
  selectDisplay(){
    let item:HTMLElement = document.querySelector("#id_"+ this.phrase.id);

    if (this.phrase.selected){
      this.phrase.selected = false;
      item.style.backgroundColor = "";
      item.style.color = "";
    } else {
      this.phrase.selected = true;
      this.setDisplayChecked();
    }
    
    return;
  }

  public disable(){
    this.disabled = true;
    let item:HTMLElement = document.querySelector("#id_"+ this.phrase.id);
    if (!this.phrase.selected){
      item.classList.add("disabled");
    }
    item.classList.add("board-disabled");
  }

  


}
