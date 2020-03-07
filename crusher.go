package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/InfamousMrS/crusher/config"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/rapidloop/skv"
)

var wsMatches map[string]*WsMatch = make(map[string]*WsMatch)
var ctxMatches map[string]string = make(map[string]string)
var sc chan os.Signal

func saveMatch(match *WsMatch) {
	store, err := skv.Open("./matches.db")
	defer store.Close()
	if err != nil {
		fmt.Print("Nothing stored can't open db")
	}
	store.Delete(match.Name)
	store.Put(match.Name, *match)
}

func loadMatch(name string) error {
	store, err := skv.Open("./matches.db")
	defer store.Close()
	if err != nil {
		fmt.Print("Nothing stored can't open db")
		return err
	}
	var match WsMatch
	err = store.Get(name, &match)
	if err != nil {
		fmt.Printf("Couldn't find the match %s\n", name)
		return err
	} else {
		wsMatches[name] = &match
		return nil
	}
}

func main() {

	// Read token from configuration
	token := config.ReadConfig()
	fmt.Printf("Token : " + token)

	// Create discordgo session
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating discord session.")
		return
	}
	// Create the Router and add the MessageCreate handl
	router := exrouter.New()

	addStartMatchRouter(router)
	addSeeMatchRouter(router)
	addLogoutRoute(router)
	addAddPlayerRouters(router)
	addListRouter(router)
	addListPlayerRouter(router)
	addStatusRouter(router)
	addUpdateRouter(router)
	addLoadMatchRouter(router)

	// shut down the bot, should be prileged.

	router.Default = router.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range router.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")

	discord.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(discord, "", discord.State.User.ID, m.Message)
	})

	// need to open the socket
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func addLoadMatchRouter(r *exrouter.Route) {
	r.On("loadmatch", func(ctx *exrouter.Context) {
		authorUsername := ctx.Msg.Author.Username
		authorDescriminator := ctx.Msg.Author.Discriminator
		authStr := authorUsername + "#" + authorDescriminator
		if authStr != "InfamousMrSatan#9232" {
			ctx.Reply(fmt.Sprintf("Sorry, %s. Only InfamousMrSatan can speak to me that way.", authorUsername))
		} else {
			if len(ctx.Args) != 2 {
				ctx.Reply("Wrong number of args, Infamous.")
			} else {
				name := ctx.Args[1]
				err := loadMatch(name)
				if err == nil {
					ctx.Reply("I think I did it!")
				} else {
					ctx.Reply(err)
				}
			}

		}
	}).Desc("Log off and shut down the bot.").Alias("signout", "exit", "quit", "fuckoff")
}

func addLogoutRoute(r *exrouter.Route) {
	r.On("signoff", func(ctx *exrouter.Context) {
		authorUsername := ctx.Msg.Author.Username
		authorDescriminator := ctx.Msg.Author.Discriminator
		authStr := authorUsername + "#" + authorDescriminator
		if authStr != "InfamousMrSatan#9232" {
			ctx.Reply(fmt.Sprintf("Sorry, %s. Only InfamousMrSatan can speak to me that way.", authorUsername))
		} else {
			ctx.Reply("Ensign Crusher, signing off.")
			sc <- syscall.SIGTERM
		}
	}).Desc("Log off and shut down the bot.").Alias("signout", "exit", "quit", "fuckoff")
}

func addSeeMatchRouter(r *exrouter.Route) {
	// set which match I'm looking at.
	r.On("seematch", func(ctx *exrouter.Context) {
		author := ctx.Msg.Author.Username
		if len(ctx.Args) == 2 {
			wsName := ctx.Args[1]
			wsMatch, exists := wsMatches[wsName]
			if exists {
				Subscribe(author, wsMatch)
				ctx.Reply(fmt.Sprintf("Aye aye, %s. Now showing you White Star [%s]", author, wsName))
			} else {
				ctx.Reply(fmt.Sprintf("Sorry Sir! White Star [%s] isn't found on any of the charts!", wsName))
			}
		} else if len(ctx.Args) == 1 {
			wsMatch, err := GetMatch(author)
			if err != nil {
				ctx.Reply(fmt.Sprintf("%s, you aren't viewing any match.", author))
			} else {
				ctx.Reply(fmt.Sprintf("%s is currently seeing match: %s", author, wsMatch.Name))
			}
		} else {
			ReplyWrongArgs(1, ctx)
		}
	}).Desc("Lets you see the stats for a certain match.").Alias("see")
}

func addStartMatchRouter(r *exrouter.Route) {
	r.On("startmatch", func(ctx *exrouter.Context) {
		if len(ctx.Args) == 2 {
			wsName := ctx.Args[1]
			_, exists := wsMatches[wsName]
			author := ctx.Msg.Author.Username
			if exists {
				// error WS already exists w/ that name
				ctx.Reply(fmt.Sprintf(
					"Sorry %s. White Star [%s] already exists!", author, wsName))
			} else {
				newMatch := NewWsMatch(wsName)
				wsMatches[wsName] = &newMatch
				Subscribe(author, &newMatch)
				ctx.Reply(fmt.Sprintf("Aye Aye, %s. I created White Star Match [%s] for you.", author, wsName))
			}
		} else {
			ReplyWrongArgs(1, ctx)
		}
	}).Desc("Creates a new mach with given name.").Alias("newmatch", "creatematch", "start", "create", "new")
}

func addAddPlayerRouters(r *exrouter.Route) {
	r.On("addfriend", func(ctx *exrouter.Context) {
		ExecuteAddPlayer(ctx, true)
	}).Desc("<name> <level> like `addfriend InfamousMrSatan 203` . no spaces.")

	r.On("addenemy", func(ctx *exrouter.Context) {
		ExecuteAddPlayer(ctx, false)
	}).Desc("<name> <level> like `addenemy BlackDeath436 225` . no spaces.")
}

func addUpdateRouter(r *exrouter.Route) {
	r.On("update", func(ctx *exrouter.Context) {
		if len(ctx.Args) == 5 {
			author := ctx.Msg.Author.Username
			match, err := GetMatch(author)
			if err != nil {
				ctx.Reply(fmt.Sprintf("Sorry %s, you aren't looking at any match!", author))
			} else {
				target := ctx.Args[1]
				targetPlayer, targetInStar := match.players[target]
				if !targetInStar {
					ctx.Reply(fmt.Sprintf("Target player [%s] isn't in this match.", target))
				} else {
					valid := true
					shiptype := ValidShipType(strings.ToLower(ctx.Args[2]))
					operation := ValidOperation(strings.ToLower(ctx.Args[3]))
					when := strings.ToLower(ctx.Args[4])
					_, err := time.ParseDuration(when)

					if shiptype == "" {
						valid = false
						ctx.Reply(
							fmt.Sprintf("You need to tell me which chip type (bs or support), but you said %s", ctx.Args[2]))
					}
					if operation == "" {
						valid = false
						ctx.Reply(
							fmt.Sprintf("You need to tell me what happened (warpin, warpout, or destroy), but you said %s", ctx.Args[3]))
					}
					if err != nil {
						valid = false
						ctx.Reply(
							fmt.Sprintf("How long ago did it happen (eg 2h30m ago) - tell me as 0h0m0s. You said: %s", ctx.Args[4]))
					}

					if valid {
						var ship *Ship
						switch shiptype {
						case "bs":
							ship = targetPlayer.Battleship
						case "support":
							ship = targetPlayer.Support
						}
						switch operation {
						case "in":
							ship.WarpedInSince(when)
						case "out":
							ship.WarpOut()
						case "destroy":
							ship.DestroyedSince(when)
						}
						ctx.Reply(fmt.Sprintf("Updated. %s", targetPlayer.Status(&systemClock)))

						saveMatch(match)
					}
				}
			}
		} else {
			ReplyWrongArgs(4, ctx)
		}
	}).Desc("<name> <BS|TS|Miner> <IN|OUT|DESTROY> <0h0m0s>")
}

func ValidOperation(operation string) string {
	switch operation {
	case "in":
		fallthrough
	case "warpin":
		fallthrough
	case "warpedin":
		return "in"
	case "out":
		fallthrough
	case "warpedout":
		fallthrough
	case "warpout":
		return "out"
	case "destroy":
		fallthrough
	case "destroyed":
		fallthrough
	case "destory":
		fallthrough
	case "destoryed":
		fallthrough
	case "kill":
		fallthrough
	case "killed":
		fallthrough
	case "killt":
		fallthrough
	case "died":
		fallthrough
	case "death":
		return "destroy"
	default:
		return ""
	}
}

func ValidShipType(shiptype string) string {
	switch shiptype {
	case "bs":
		fallthrough
	case "battleship":
		return "bs"
	case "ts":
		fallthrough
	case "transport":
		fallthrough
	case "miner":
		fallthrough
	case "support":
		fallthrough
	case "utility":
		return "support"
	default:
		return ""
	}
}

func addStatusRouter(r *exrouter.Route) {
	r.On("status", func(ctx *exrouter.Context) {
		author := ctx.Msg.Author.Username
		match, err := GetMatch(author)
		if err != nil {
			ctx.Reply(fmt.Sprintf("Sorry %s, you aren't looking at any match!", author))
		} else {
			reply := fmt.Sprintf("```\nMatch Status [%s]\n\n", match.Name)
			reply += "Friendly Status:\n"
			for _, p := range match.Friendlies {
				reply += fmt.Sprintf("%s\n", Status(&p))
			}
			reply += "\nEnemy Status:\n"
			for _, p := range match.Enemies {
				reply += fmt.Sprintf("%s\n", Status(&p))
			}
			reply += fmt.Sprintf("```\n")
			ctx.Reply(reply)
		}

	}).Desc("Print the status of the match.")
}

var systemClock = SystemClock{}

func Status(player *Player) string {
	return player.Status(&systemClock)
}

func addListPlayerRouter(r *exrouter.Route) {
	r.On("players", func(ctx *exrouter.Context) {
		author := ctx.Msg.Author.Username
		match, err := GetMatch(author)
		if err != nil {
			ctx.Reply(fmt.Sprintf("Sorry %s, you aren't looking at any match!", author))
		} else {
			reply := "These are the players in match [" + match.Name + "]\n```"
			reply += "Friendlies:\n"
			for _, p := range match.Friendlies {
				reply += fmt.Sprintf("%s\n", p.Name)
			}
			reply += "\nEnemies:\n"
			for _, p := range match.Enemies {
				reply += fmt.Sprintf("%s\n", p.Name)
			}
			reply += "```"
			ctx.Reply(reply)
		}
	}).Desc("list the matches.").Alias("listplayers")
}

func addListRouter(r *exrouter.Route) {
	r.On("list", func(ctx *exrouter.Context) {
		author := ctx.Msg.Author.Username
		match, err := GetMatch(author)
		authormatch := ""
		if err == nil {
			authormatch = match.Name
		}
		reply := "Current White Star Matches.```\n"
		for k, _ := range wsMatches {
			if k == authormatch {
				reply += "*"
			}
			reply += k + "\n"
		}
		reply += fmt.Sprintf("```\n Oh and by the way - you're currently looking at the %s", authormatch)
		ctx.Reply(reply)
	}).Desc("list the matches.").Alias("listmatches")
}

func GetMatch(userId string) (*WsMatch, error) {
	matchName, ctxExists := ctxMatches[userId]
	if ctxExists {
		match, matchExists := wsMatches[matchName]
		if matchExists {
			return match, nil
		}
		return nil, fmt.Errorf("Couldn't find match %s", matchName)
	}
	return nil, fmt.Errorf("Player %s is not subscribed to any match", userId)
}

// Subscribe to a match
func Subscribe(userId string, wsMatch *WsMatch) error {
	_, exists := wsMatches[wsMatch.Name]
	if exists {
		ctxMatches[userId] = wsMatch.Name
		return nil
	}
	return errors.New(fmt.Sprintf("No such WS Match [%s]", wsMatch.Name))
}

// ReplyWrongArgs helper
func ReplyWrongArgs(requiredInputs int, ctx *exrouter.Context) {
	author := ctx.Msg.Author.Username
	ctx.Reply(fmt.Sprintf("Sorry, %s. The command %s requires %d inputs. But I got %d!",
		author, ctx.Args[0], requiredInputs, len(ctx.Args)-1))
}

func ExecuteAddPlayer(ctx *exrouter.Context, friendly bool) {
	if len(ctx.Args) == 3 {
		author := ctx.Msg.Author.Username
		playername := ctx.Args[1]
		playerlevel, err := strconv.Atoi(ctx.Args[2])
		if err != nil {
			ctx.Reply(fmt.Sprintf("%s is a funny number for %s's level.", ctx.Args[2], playername))
			return
		}
		wsMatch, err := GetMatch(author)
		if err == nil {
			player := NewPlayer(playername, playerlevel)
			wsMatch.addPlayer(player, friendly)
			teamString := "friendly"
			if !friendly {
				teamString = "enemy"
			}
			ctx.Reply(fmt.Sprintf("Aye Aye, Sir. Player %s added to the %s roster!", playername, teamString))
		} else {
			ctx.Reply(fmt.Sprintf("Sorry Sir! I can't do that because %s", err))
		}
	} else {
		ReplyWrongArgs(2, ctx)
	}
}
