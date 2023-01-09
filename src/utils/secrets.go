package utils

import (
	"fmt"
	"log"
	"r00t2.io/gosecret"
)

func SetSecret(collectionName string, secretName string, secretValue string, itemAttrs map[string]string) error {
	var err error
	var service *gosecret.Service
	var collection *gosecret.Collection
	var secret *gosecret.Secret

	if service, err = gosecret.NewService(); err != nil {
		return err
	}
	defer func(service *gosecret.Service) {
		err := service.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(service)

	if collection, err = service.GetCollection(collectionName); err != nil {
		if err == gosecret.ErrDoesNotExist {
			collection, err = service.CreateCollection(collectionName)
			if err != nil {
				return err
			}
		}
		return err
	}

	err = collection.Unlock()
	if err != nil {
		return err
	}

	if secretValue == "" {
		panic("secretValue is empty")
	}

	items, err := collection.Items()
	if err != nil {
		return err
	}

	var item *gosecret.Item
	for _, i := range items {
		label, err := i.Label()
		if err != nil {
			return err
		}

		if label == secretName {
			item = i
			break
		}
	}

	if item != nil {
		err := item.Delete()
		if err != nil {
			return err
		}
	}

	secret = gosecret.NewSecret(
		service.Session,
		[]byte{},
		[]byte(secretValue),
		"text/plain",
	)

	if _, err = collection.CreateItem(
		secretName,
		itemAttrs,
		secret,
		true,
	); err != nil {
		return err
	}

	return nil
}

func GetSecret(collectionName string, secretName string) (string, error) {
	var err error
	var service *gosecret.Service
	var collection *gosecret.Collection

	// All interactions with SecretService start with initiating a Service connection.
	if service, err = gosecret.NewService(); err != nil {
		return "", err
	}
	defer func(service *gosecret.Service) {
		err := service.Close()
		if err != nil {
			log.Panicln(err)
			return
		}
	}(service)

	// And unless operating directly on a Service via its methods, you probably need a Collection as well.
	if collection, err = service.GetCollection(collectionName); err != nil {
		if err == gosecret.ErrDoesNotExist {
			return "", nil
		}

		return "", err
	}

	err = collection.Unlock()
	if err != nil {
		return "", err
	}

	var itemResults []*gosecret.Item

	if itemResults, err = collection.Items(); err != nil {
		return "", err
	}

	var itemLabel string
	for _, item := range itemResults {
		if itemLabel, err = item.Label(); err != nil {
			fmt.Printf("Can't read label for item at path '%v'\n", item.Dbus.Path())
			continue
		}

		if itemLabel != secretName {
			continue
		}

		err := item.Unlock()
		if err != nil {
			return "", err
		}

		return string(item.Secret.Value), nil
	}

	return "", ErrSecretNotFound
}
