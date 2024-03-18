package main

import (
	testSuite "cloud-render/cmd/tests/intergration/suite"
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/db/redis"
	"cloud-render/internal/lib/config"
	"log"
	"os"
	"strconv"
)

func main() {
	all := 0     // counter of all tests
	passed := 0  // counter of passed tests
	failed := 0  // counter of failed tests
	errored := 0 // counter of errored tests
	_ = all

	authCfgPath := os.Getenv("AUTH_CONFIG_PATH")
	apiCfgPath := os.Getenv("API_CONFIG_PATH")
	bufferCfgPath := os.Getenv("BUFFER_CONFIG_PATH")

	log.Println(authCfgPath)
	log.Println(apiCfgPath)
	log.Println(bufferCfgPath)

	authCfg := config.MustLoad(authCfgPath)
	apiCfg := config.MustLoad(apiCfgPath)
	bufferCfg := config.MustLoad(bufferCfgPath)

	authPG, err := postgres.New(authCfg.DB)
	if err != nil {
		log.Fatalf("failed to initialize auth storage: %s", err.Error())
		os.Exit(-1)
	}
	defer authPG.Close()

	postgres.MigrateTop(authPG, authCfg.DB.MigrationsPath)
	defer postgres.DropMigrations(authPG, authCfg.DB.MigrationsPath)

	apiPG, err := postgres.New(apiCfg.DB)
	if err != nil {
		log.Fatalf("failed to initialize apit storage: %s", err.Error())
		os.Exit(-1)
	}
	defer apiPG.Close()

	postgres.MigrateTop(apiPG, apiCfg.DB.MigrationsPath)
	defer postgres.DropMigrations(apiPG, apiCfg.DB.MigrationsPath)

	// Redis
	client, err := redis.New(apiCfg)
	if err != nil {
		log.Fatalf("failed to initialize redis: %s", err.Error())
		os.Exit(-1)
	}
	defer client.Close()

	suite := testSuite.SetupSuite(authCfg, apiCfg, bufferCfg)

	// TESTS

	// SendRequest
	redis.Clear(client, bufferCfg.Redis.QueueName)
	isPassed, msg, err := suite.TestSendRequest()
	all++
	if err != nil {
		errored++
		log.Printf("TestSendRequest [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestSendRequest [PASS]")
	} else {
		failed++
		log.Printf("TestSendRequest [FAIL]: %s\n", msg)
	}

	// UserEdit
	isPassed, msg, err = suite.TestUserEdit()
	all++
	if err != nil {
		errored++
		log.Printf("TestUserEdit [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestUserEdit [PASS]")
	} else {
		failed++
		log.Printf("TestUserEdit [FAIL]: %s\n", msg)
	}

	// SubscribeUser
	isPassed, msg, err = suite.TestSubscribeUser()
	all++
	if err != nil {
		errored++
		log.Printf("TestSubscribeUser [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestSubscribeUser [PASS]")
	} else {
		failed++
		log.Printf("TestSubscribeUser [FAIL]: %s\n", msg)
	}

	// SignInEdit
	isPassed, msg, err = suite.TestEditSignIn()
	all++
	if err != nil {
		errored++
		log.Printf("TestEditSignIn [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestEditSignIn [PASS]")
	} else {
		failed++
		log.Printf("TestEditSignIn [FAIL]: %s\n", msg)
	}

	// OrdersSend
	isPassed, msg, err = suite.TestOrdersSend()
	all++
	if err != nil {
		errored++
		log.Printf("TestOrdersSend [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestOrdersSend [PASS]")
	} else {
		failed++
		log.Printf("TestOrdersSend [FAIL]: %s\n", msg)
	}

	// OrdersUpdate
	isPassed, msg, err = suite.TestOrdersUpdate(strconv.Itoa(all + 1))
	all++
	if err != nil {
		errored++
		log.Printf("TestOrdersUpdate [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestOrdersUpdate [PASS]")
	} else {
		failed++
		log.Printf("TestOrdersUpdate [FAIL]: %s\n", msg)
	}

	// OrdersDelete
	isPassed, msg, err = suite.TestOrdersDelete()
	all++
	if err != nil {
		errored++
		log.Printf("TestOrdersDelete [ERROR]: %s\n", err.Error())
	} else if isPassed {
		passed++
		log.Printf("TestOrdersDelete [PASS]")
	} else {
		failed++
		log.Printf("TestOrdersDelete [FAIL]: %s\n", msg)
	}

	log.Printf("Tests summary:\n\tPASSED [%d/%d]\n\tFAILED [%d/%d]\n\tERRORED [%d/%d]\n", passed, all, failed, all, errored, all)
	if passed == all {
		log.Println("TESTS PASSED")
	} else {
		log.Println("TESTS FAILED")
	}
}
