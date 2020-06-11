import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BoardComponent } from './view/board/board.component';
import { ItemComponent } from './view/board/item/item.component';
import { HttpClientModule }    from '@angular/common/http';
import { environment } from "src/environments/environment";
import { AngularFireModule } from "@angular/fire";
import { AngularFirestoreModule } from "@angular/fire/firestore";
import { LoginComponent } from './view/login/login.component';
import { MessagesComponent } from './view/messages/messages.component';
import { ReplacePipe } from './replace.pipe';
import { AdminComponent } from './view/admin/admin.component';
import { ToolbarComponent } from './view/toolbar/toolbar.component';
import { ManageBoardsComponent } from './view/manage/manage-boards/manage-boards.component';
import { ManageMasterComponent } from './view/manage/manage-master/manage-master.component';
import { AdminPhrasesComponent } from './view/admin/admin-phrases/admin-phrases.component';
import { AdminPhraseFormComponent } from './view/admin/admin-phrases/admin-phrase-form/admin-phrase-form.component';
import { AdminGameComponent } from './view/admin/admin-game/admin-game.component';
import { GamepickerComponent } from './view/login/gamepicker/gamepicker.component';
import { ManageComponent } from './view/manage/manage.component';
import { ManagePhrasesComponent } from './view/manage/manage-phrases/manage-phrases.component';
import { ManagePhraseFormComponent } from './view/manage/manage-phrases/manage-phrase-form/manage-phrase-form.component';
import { ManageAdminsComponent } from './view/manage/manage-admins/manage-admins.component';
import { AdminAdminsComponent } from './view/admin/admin-admins/admin-admins.component';
import { FooterComponent } from './view/footer/footer.component';
import { AdminGamesComponent } from './view/admin/admin-games/admin-games.component';
import { ProgressbarComponent } from './view/widgets/progressbar/progressbar.component';

@NgModule({
  declarations: [
    AppComponent,
    BoardComponent,
    ItemComponent,
    LoginComponent,
    ReplacePipe,
    MessagesComponent,
    AdminComponent,
    ToolbarComponent,
    ManageBoardsComponent,
    ManageMasterComponent,
    AdminPhrasesComponent,
    AdminPhraseFormComponent,
    AdminGameComponent,
    GamepickerComponent,
    ManageComponent,
    ManagePhrasesComponent,
    ManagePhraseFormComponent,
    ManageAdminsComponent,
    AdminAdminsComponent,
    FooterComponent,
    AdminGamesComponent,
    ProgressbarComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    AppRoutingModule,
    AngularFireModule.initializeApp(environment.firebaseConfig),
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
