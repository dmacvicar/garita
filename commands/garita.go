package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dmacvicar/garita/api"
)

var GaritaCmd = &cobra.Command{
	Use:   "garita",
	Short: "light v2 auth server for docker",
	Long:  "starts the garita server",
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig()
		server()
	},
}

var garitaCmdV *cobra.Command

// cmd line flags
var insecureHttpF bool
var listenPortF int
var htpasswdPathF string
var keyPathF string
var tlsCertPathF string
var tlsKeyPathF string
var configPathF string

//Execute adds all child commands to the root command GaritaCmd and sets flags appropriately.
func Execute() {
	// add the subcommands here in the future if any
	GaritaCmd.Execute()
}

func init() {
	GaritaCmd.Flags().StringVarP(&keyPathF, "key", "k", "./server.key", "Auth token secret key")
	GaritaCmd.Flags().BoolVarP(&insecureHttpF, "http", "x", false, "use plain HTTP")
	GaritaCmd.Flags().IntVarP(&listenPortF, "port", "p", 443, "Port to listen to")
	GaritaCmd.Flags().StringVarP(&htpasswdPathF, "htpasswd", "w", "./htpasswd", "htpasswd file")
	GaritaCmd.Flags().StringVarP(&tlsCertPathF, "tlscert", "s", "./server.crt", "TLS certificate")
	GaritaCmd.Flags().StringVarP(&tlsKeyPathF, "tlskey", "y", "./server.key", "TLS secret key")

	GaritaCmd.Flags().StringVarP(&configPathF, "config", "c", "", "Configuration file. Command line options override settings in the configuration file")

	garitaCmdV = GaritaCmd
}

func initializeConfig() {

	if garitaCmdV.Flags().Lookup("config").Changed {
		viper.SetConfigFile(configPathF)
		err := viper.ReadInConfig()
		if err != nil {
			log.Println("Unable to locate configuration file: " + configPathF)
		} else {
			log.Println(fmt.Sprintf("Using configuration %s", viper.ConfigFileUsed()))
		}
	}

	viper.SetDefault("key", "server.key")
	viper.SetDefault("port", 443)
	viper.SetDefault("http", false)
	viper.SetDefault("htpasswd", "htpasswd")
	viper.SetDefault("tlscert", "server.crt")
	viper.SetDefault("tlskey", "server.key")

	if garitaCmdV.Flags().Lookup("key").Changed {
		viper.Set("key", &keyPathF)
	}

	if garitaCmdV.Flags().Lookup("http").Changed {
		viper.Set("http", &insecureHttpF)
	}

	if garitaCmdV.Flags().Lookup("port").Changed {
		viper.Set("port", &listenPortF)
	}

	if garitaCmdV.Flags().Lookup("htpasswd").Changed {
		viper.Set("htpasswd", &htpasswdPathF)
	}

	if garitaCmdV.Flags().Lookup("tlskey").Changed {
		viper.Set("tlskey", &tlsKeyPathF)
	}

	if garitaCmdV.Flags().Lookup("tlscert").Changed {
		viper.Set("tlscert", &tlsCertPathF)
	}
}

func server() {

	insecureHttp := viper.GetBool("http")
	listenPort := viper.GetInt("port")
	htpasswdPath := viper.GetString("htpasswd")
	keyPath := viper.GetString("key")
	tlsCertPath := viper.GetString("tlscert")
	tlsKeyPath := viper.GetString("tlskey")

	if _, err := os.Stat(htpasswdPath); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", htpasswdPath)
		return
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", keyPath)
		return
	}

	// tls requires both cert and key
	if !insecureHttp {
		if _, err := os.Stat(tlsCertPath); os.IsNotExist(err) {
			fmt.Printf("no such file or directory: %s", tlsCertPath)
			return
		}

		if _, err := os.Stat(tlsKeyPath); os.IsNotExist(err) {
			fmt.Printf("no such file or directory: %s", tlsKeyPath)
			return
		}
	}

	tokenHandler := api.NewGaritaTokenHandler(htpasswdPath, keyPath)

	router := mux.NewRouter()
	router.Handle("/v2/token", tokenHandler)
	log.Printf("Listening...:%d", listenPort)

	if insecureHttp {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), router))
	} else {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", listenPort), tlsCertPath, tlsKeyPath, router))
	}
}
