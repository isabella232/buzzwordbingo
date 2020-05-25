package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// I stole this code from firestore/collref.go basically it generates the ids
// so I can use batch sets instead of adds for anything
const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func uniqueID() string {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("agent: crypto/rand.Read error: %v", err))
	}
	for i, byt := range b {
		b[i] = alphanum[int(byt)%len(alphanum)]
	}
	return string(b)
}

// Agent is a go between for the main application and firestore.
type Agent struct {
	ProjectID string
}

func (a *Agent) log(msg string) {
	if noisy {
		log.Printf("Firestore: %s\n", msg)
	}
}

func (a *Agent) getClient() (*firestore.Client, error) {
	if client != nil {
		return client, nil
	}
	a.log("Getting New Client")
	return firestore.NewClient(context.Background(), a.ProjectID)
}

// IsAdmin tests if a give player is in the admin group by email
func (a *Agent) IsAdmin(email string) (bool, error) {
	var err error
	client, err = a.getClient()
	if err != nil {
		return false, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("See if user exists")
	doc, err := client.Collection("admins").Doc(email).Get(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "code = NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to get game: %v", err)
	}

	return doc.Exists(), nil

}

// GetPhrases fetches the master list of Phrases for populating Games
func (a *Agent) GetPhrases() ([]Phrase, error) {

	p := []Phrase{}

	var err error
	client, err = a.getClient()
	if err != nil {
		return p, fmt.Errorf("Failed to create client: %v", err)
	}

	a.log("Getting Phrases")
	iter := client.Collection("phrases").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return p, fmt.Errorf("Failed to iterate: %v", err)
		}
		var phrase Phrase
		dataMap := doc.Data()

		phrase.ID = dataMap["id"].(string)
		phrase.Text = dataMap["text"].(string)

		p = append(p, phrase)
	}

	return p, nil
}

// NewGame will create a new game in the database and initialize it.
func (a *Agent) NewGame(name string, p Player) (Game, error) {
	g := Game{}

	phrases, err := a.GetPhrases()
	if err != nil {
		return g, fmt.Errorf("failed to get phrases: %v", err)
	}

	g.ID = uniqueID()
	g.Admins = append(g.Admins, p)
	g.Name = name
	g.Active = true
	g.master.Load(phrases)

	client, err = a.getClient()

	if err != nil {
		return g, fmt.Errorf("failed to create client: %v", err)
	}

	batch := client.Batch()
	a.log("Creating new game")

	gref := client.Collection("games").Doc(g.ID)
	batch.Set(gref, g)

	a.log("Adding phrases to new game")
	for _, v := range g.master.Records {
		ref := client.Collection("games").Doc(g.ID).Collection("records").Doc(v.Phrase.ID)
		batch.Set(ref, v)
	}

	aref := client.Collection("games").Doc(g.ID).Collection("admins").Doc(p.Email)
	batch.Set(aref, p)

	m := Message{}
	m.SetText("Game has begun!")
	m.SetAudience("all")

	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	mref := client.Collection("games").Doc(g.ID).Collection("messages").Doc(timestamp)
	batch.Set(mref, m)

	_, err = batch.Commit(ctx)
	if err != nil {
		return g, fmt.Errorf("failed to add records to database: %v", err)
	}

	return g, nil
}

// GetGame gets a given game from the database
func (a *Agent) GetGame(id string) (Game, error) {
	g := Game{}
	var err error
	client, err = a.getClient()
	if err != nil {
		return g, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("Getting existing game")
	doc, err := client.Collection("games").Doc(id).Get(ctx)
	if err != nil {
		return g, fmt.Errorf("failed to get game: %v", err)
	}

	doc.DataTo(&g)
	g.ID = id
	g, err = a.loadGameWithRecords(g)
	if err != nil {
		return g, fmt.Errorf("failed to load records for game: %v", err)
	}

	g, err = a.loadGameWithPlayers(g)
	if err != nil {
		return g, fmt.Errorf("failed to load players for game: %v", err)
	}

	g, err = a.loadGameWithAdmins(g)
	if err != nil {
		return g, fmt.Errorf("failed to load admins for game: %v", err)
	}

	return g, nil
}

func (a *Agent) loadGameWithRecords(g Game) (Game, error) {
	var err error
	client, err = a.getClient()
	if err != nil {
		return g, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("Loading records from game")
	iter := client.Collection("games").Doc(g.ID).Collection("records").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return g, fmt.Errorf("failed getting game records: %v", err)
		}
		r := Record{}
		doc.DataTo(&r)
		g.master.Records = append(g.master.Records, r)
	}

	return g, nil
}

func (a *Agent) loadGameWithPlayers(g Game) (Game, error) {
	var err error
	client, err = a.getClient()
	if err != nil {
		return g, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("Loading players from game")
	iter := client.Collection("games").Doc(g.ID).Collection("players").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return g, fmt.Errorf("failed getting game records: %v", err)
		}
		p := Player{}
		doc.DataTo(&p)
		g.Players = append(g.Players, p)
	}

	return g, nil
}

func (a *Agent) loadGameWithAdmins(g Game) (Game, error) {
	var err error
	client, err = a.getClient()
	if err != nil {
		return g, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("Loading players from game")
	iter := client.Collection("games").Doc(g.ID).Collection("admins").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return g, fmt.Errorf("failed getting game records: %v", err)
		}
		p := Player{}
		doc.DataTo(&p)
		g.Admins = append(g.Admins, p)
	}

	return g, nil
}

// AddMessagesToGame broadcasts a message to the game players
func (a *Agent) AddMessagesToGame(g Game, m []Message) error {

	var err error
	client, err = a.getClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}

	batch := client.Batch()
	for _, v := range m {
		a.log("Adding message to game")
		timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		v.ID = timestamp
		ref := client.Collection("games").Doc(g.ID).Collection("messages").Doc(timestamp)
		batch.Set(ref, v)
	}

	_, err = batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to send messages : %v", err)
	}

	return nil
}

// SaveGame records a game to firestore.
func (a *Agent) SaveGame(g Game) error {

	var err error
	client, err = a.getClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}

	a.log("Save game")
	batch := client.Batch()
	gref := client.Collection("games").Doc(g.ID)
	batch.Set(gref, g)

	for _, v := range g.Players {
		ref := client.Collection("games").Doc(g.ID).Collection("players").Doc(v.Email)
		batch.Set(ref, v)
	}

	for _, v := range g.Admins {
		ref := client.Collection("games").Doc(g.ID).Collection("admins").Doc(v.Email)
		batch.Set(ref, v)
	}

	_, err = batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to save game to database: %v", err)
	}

	return nil
}

// GetBoardForPlayer returns the board for a given player
func (a *Agent) GetBoardForPlayer(id string, p Player) (Board, error) {
	b := Board{}
	var err error
	client, err = a.getClient()
	if err != nil {
		return b, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("get board for player")
	iter := client.Collection("games").Doc(id).Collection("boards").Where("game", "==", id).Where("player.email", "==", p.Email).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return b, fmt.Errorf("failed to iterate over b from firestore: %v", err)
		}
		doc.DataTo(&b)
		b.ID = doc.Ref.ID
		break
	}

	if b.ID != "" {
		b, err = a.loadBoardWithPhrases(b)
		if err != nil {
			return b, fmt.Errorf("failed to load phrases for board: %v", err)
		}
	}

	return b, nil

}

// GetBoard retrieves a specifc board from firestore
func (a *Agent) GetBoard(bid, gid string) (Board, error) {
	b := Board{}
	var err error
	client, err = a.getClient()
	if err != nil {
		return b, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("Getting board")
	doc, err := client.Collection("games").Doc(gid).Collection("boards").Doc(bid).Get(ctx)
	if err != nil {
		return b, fmt.Errorf("failed to get board: %v", err)
	}

	doc.DataTo(&b)
	b.ID = bid
	b, err = a.loadBoardWithPhrases(b)
	if err != nil {
		return b, fmt.Errorf("failed to load phrases for board: %v", err)
	}
	return b, nil
}

func (a *Agent) loadBoardWithPhrases(b Board) (Board, error) {
	var err error
	client, err = a.getClient()
	if err != nil {
		return b, fmt.Errorf("failed to create client: %v", err)
	}

	a.log("Adding phrases to existing board")
	iter := client.Collection("games").Doc(b.Game).Collection("boards").Doc(b.ID).Collection("phrases").OrderBy("displayorder", firestore.Asc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return b, fmt.Errorf("failed getting board records: %v", err)
		}
		p := Phrase{}
		doc.DataTo(&p)
		b.Phrases = append(b.Phrases, p)
	}

	return b, nil
}

// DeleteBoard delete a specifc board from firestore
func (a *Agent) DeleteBoard(bid, gid string) error {
	var err error
	client, err = a.getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	batch := client.Batch()
	a.log("Deleting board")
	bref := client.Collection("games").Doc(gid).Collection("boards").Doc(bid)
	batch.Delete(bref)
	a.log("removing phrases from board")
	ref := client.Collection("games").Doc(gid).Collection("boards").Doc(bid).Collection("phrases")
	for {
		// Get a batch of documents
		iter := ref.Limit(100).Documents(ctx)

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to clean phrases from firestore: %v", err)
			}

			a.log(fmt.Sprintf("removing phrase %s from board", doc.Ref.ID))
			batch.Delete(doc.Ref)
		}
		_, err = batch.Commit(ctx)
		if err != nil {
			return fmt.Errorf("failed to clean messages from firestore: %v", err)
		}

		return nil

	}

}

// SaveBoard persists a board to firestore
func (a *Agent) SaveBoard(b Board) (Board, error) {
	var err error
	client, err = a.getClient()
	if err != nil {
		return b, fmt.Errorf("failed to create client: %v", err)
	}

	if b.ID == "" {
		b.ID = uniqueID()
	}

	a.log("Starting batch operation")
	batch := client.Batch()
	bref := client.Collection("games").Doc(b.Game).Collection("boards").Doc(b.ID)
	batch.Set(bref, b)

	pref := client.Collection("games").Doc(b.Game).Collection("players").Doc(b.Player.Email)
	batch.Set(pref, b.Player)

	for _, v := range b.Phrases {
		ref := client.Collection("games").Doc(b.Game).Collection("boards").Doc(b.ID).Collection("phrases").Doc(v.ID)
		batch.Set(ref, v)
	}

	_, err = batch.Commit(ctx)
	if err != nil {
		return b, fmt.Errorf("failed to add records to database: %v", err)
	}

	return b, nil

}

// UpdatePhrase records clicks on the board and the game
func (a *Agent) UpdatePhrase(b Board, p Phrase, r Record) error {
	var err error
	client, err = a.getClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}

	a.log("Starting batch operation")
	batch := client.Batch()

	a.log("Updating phrase on board")
	bref := client.Collection("games").Doc(b.Game).Collection("boards").Doc(b.ID).Collection("phrases").Doc(p.ID)
	batch.Set(bref, p)

	a.log("Updating game record")
	gref := client.Collection("games").Doc(b.Game).Collection("records").Doc(r.Phrase.ID)
	batch.Set(gref, r)

	a.log("Updating board to bingo")
	bingoref := client.Collection("games").Doc(b.Game).Collection("boards").Doc(b.ID)
	update := map[string]interface{}{"bingodeclared": b.BingoDeclared}
	batch.Set(bingoref, update, firestore.MergeAll)

	a.log("Committing Batch")
	_, err = batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to update phrase: %v", err)
	}

	return nil
}

// GetGamesForPlayer fetches the list of all games a user in currently in.
func (a *Agent) GetGamesForPlayer(email string) (Games, error) {

	g := []Game{}

	var err error
	client, err = a.getClient()
	if err != nil {
		return g, fmt.Errorf("Failed to create client: %v", err)
	}

	refs := []*firestore.DocumentRef{}
	a.log("Getting Boards for player")
	iter := client.CollectionGroup("players").Where("email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return g, fmt.Errorf("Failed to iterate: %v", err)
		}

		refs = append(refs, doc.Ref.Parent.Parent)
	}

	a.log("Getting Games for player")
	snapshots, err := client.GetAll(ctx, refs)
	if err != nil {
		return g, fmt.Errorf("Failed to get games: %v", err)
	}

	for _, v := range snapshots {
		game := Game{}
		v.DataTo(&game)
		game.ID = v.Ref.ID
		g = append(g, game)
	}

	return g, nil
}
