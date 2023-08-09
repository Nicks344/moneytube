package mutations

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/accounts"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getUserAgent() string {
	content, err := ioutil.ReadFile("user-agents.txt")
	if err != nil {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	}

	lines := strings.Split(string(content), "\n")
	var agents []string
	for _, line := range lines {
		line = strings.Trim(line, " \r")
		if line != "" {
			agents = append(agents, line)
		}
	}

	if len(agents) == 0 {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	}

	return agents[rand.Intn(len(agents))]
}

func addOrEditAccount() *graphql.Field {
	return &graphql.Field{
		Args:        gqlstructs.AccountInput,
		Type:        gqlstructs.AccountOutput,
		Description: "Add new or edit account",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var account moneytubemodel.Account

			mapstructure.Decode(params.Args["account"], &account)

			if account.ID > 0 {
				old, err := model.GetAccount(account.ID)
				if err != nil {
					return nil, err
				}

				old.Proxy = account.Proxy
				old.Group = account.Group
				account = old
			} else {
				account.UserAgent = getUserAgent()
			}

			if params.Args["cookieFile"] != nil {
				cookieFile := params.Args["cookieFile"].(string)

				if cookieFile != "" {
					if err := accounts.ImportCookies(account, cookieFile); err != nil {
						return nil, err
					}
				}
			}

			if err := model.SaveAccount(&account); err != nil {
				return nil, err
			}

			return account, nil
		},
	}
}

func deleteAccount() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Delete account",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			accID := params.Args["id"].(int)
			err := model.DeleteAccount(accID)
			if err != nil {
				return err == nil, err
			}
			err = model.DeleteUploadTasksWithAccountID(accID)
			return err == nil, err
		},
	}
}

func updateAccount() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Update account",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			account, err := model.GetAccount(params.Args["id"].(int))
			if err != nil {
				return false, err
			}
			account.Status = moneytubemodel.ASUpdate
			events.OnAccountUpdated(account)
			go func() {
				acc, err := accounts.GetInfo(account)
				if err != nil {
					account.Status = moneytubemodel.ASError
					account.ErrorMessage = err.Error()
				} else {
					account = acc
					account.Status = moneytubemodel.ASReady
				}
				model.SaveAccount(&account)
				events.OnAccountUpdated(account)
			}()
			return true, nil
		},
	}
}

func updateAllAccounts() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Delete account",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			accs, err := model.GetAccountsByStatus(moneytubemodel.ASReady)
			if err != nil {
				return false, err
			}
			for _, account := range accs {
				account.Status = moneytubemodel.ASUpdate
				events.OnAccountUpdated(account)
				go func(account moneytubemodel.Account) {
					acc, err := accounts.GetInfo(account)
					if err != nil {
						account.Status = moneytubemodel.ASError
						account.ErrorMessage = err.Error()
					} else {
						account = acc
						account.Status = moneytubemodel.ASReady
					}
					model.SaveAccount(&account)
					events.OnAccountUpdated(account)
				}(account)
			}

			return true, nil
		},
	}
}

func openAccountBrowser() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Open account's browser",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			account, err := model.GetAccount(params.Args["id"].(int))
			if err != nil {
				return false, err
			}
			err = accounts.OpenAndShow(account)
			return err == nil, err
		},
	}
}

func deleteGroup() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"group": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Delete group",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			err := model.DeleteGroup(params.Args["group"].(string))
			return err == nil, err
		},
	}
}

func importCookies() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"file": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Import account cookies",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			acc, err := model.GetAccount(p.Args["id"].(int))
			if err != nil {
				return nil, err
			}

			file := p.Args["file"].(string)

			err = accounts.ImportCookies(acc, file)
			return err == nil, err
		},
	}
}
