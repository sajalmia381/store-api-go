package utils

import (
	"math/rand"
	"time"

	"github.com/gosimple/slug"
	"github.com/sajalmia381/store-api/src/v1/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var smallChars = "abcdefghijklmnopqrstuvwxyz"

func GenerateUniqueSlug(title string, collectionName string, skip_slugs ...string) string {
	newSlug := slug.MakeLang(title, "en")
	collection := db.GetDmManager().DB.Collection(collectionName)
	for {
		result := collection.FindOne(db.GetDmManager().Ctx, bson.M{"slug": newSlug})
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
		} else {
			for _, s := range skip_slugs {
				if s == newSlug {
					return newSlug
				}
			}
			newSlug = newSlug + "-" + GenerateRandomString(5, &smallChars)
		}
	}
	return newSlug
}

func GenerateFakeUniqueSlug(title string, withPrefix bool, skip_slugs ...string) string {
	// For Fake generator
	newSlug := slug.MakeLang(title, "en")
	if withPrefix {
		newSlug = newSlug + "-" + GenerateRandomString(5, &smallChars)
	}
	return newSlug
}

// Generate Random String
func mergeCharSets(charset ...*string) string {
	_charset := ""
	for _, char := range charset {
		if char == nil {
			continue
		}
		if char != nil {
			_charset = _charset + *char
		}
	}
	return _charset
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateRandomString(length int16, charSets ...*string) string {
	_charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if len(charSets) > 0 {
		_charset = mergeCharSets(charSets...)
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = _charset[seededRand.Intn(len(_charset))]
	}
	return string(b)
}
