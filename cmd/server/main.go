package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/auth"
	"github.com/kevinwylder/sbvision/dynamo"
	"github.com/kevinwylder/sbvision/video/encoder"

	"github.com/gorilla/websocket"
)

type serverContext struct {
	upgrader   websocket.Upgrader
	auth       *auth.JWTVerifier
	videoCache map[int64]*sbvision.Video
	remotes    map[string]*remoteSession
	processes  *encoder.VideoRequestManager
	ddb        *dynamo.SBDatabase
	batch      *batch.Batch
}

func main() {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ddb, err := dynamo.FindTables(session)
	if err != nil {
		log.Fatal(err)
	}

	server := &serverContext{
		ddb:   ddb,
		batch: batch.New(session),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 20 * 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		remotes:   make(map[string]*remoteSession),
		auth:      auth.NewJWTVerifier(ddb, "https://cognito-idp.us-west-2.amazonaws.com/us-west-2_dHWlJDm4T/.well-known/jwks.json"),
		processes: encoder.NewVideoRequestManager(session),
	}

	fmt.Println("Starting server")
	err = http.ListenAndServe(":"+os.Getenv("PORT"), server)
	fmt.Println(err)
}
