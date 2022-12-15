package main

import (
	"at.ourproject/energystore/calculation"
	"at.ourproject/energystore/mqttclient"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"at.ourproject/energystore/config"
	"at.ourproject/energystore/graph"
	"at.ourproject/energystore/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/spf13/viper"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var configPath = flag.String("configPath", ".", "Configfile Path")
	flag.Parse()

	println("-> \nRead Config")
	config.ReadConfig(*configPath)

	ctx, _ := context.WithCancel(context.Background())
	SetupMqttDispatcher(ctx)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func SetupMqttDispatcher(ctx context.Context) {
	streamer, err := mqttclient.NewMqttStreamer()
	if err != nil {
		panic(err)
	}

	worker := map[string]mqttclient.Executor{}
	energyTopicPrefix := viper.GetString("mqtt.energySubscriptionTopic")
	communityIds := viper.GetStringSlice("mqtt.communityIds")
	for _, id := range communityIds {
		worker[fmt.Sprintf("%s/%s", energyTopicPrefix, id)] = calculation.NewMqttEnergyImporter(id)
	}
	//worker := map[string]mqttclient.Executor{
	//	"eda/response/energy/rc100130": calculation.NewMqttEnergyImporter("rc100130"),
	//}

	dispatcher := mqttclient.NewDispatcher(ctx, streamer, worker)
	_ = dispatcher
}
